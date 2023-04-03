package otlpinf

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	yson "github.com/ghodss/yaml"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leoparente/otlpinf/config"
	"github.com/leoparente/otlpinf/runner"
	"gopkg.in/yaml.v3"
)

func messageReturn(value string) []byte {
	ret := strings.Join([]string{`{"message":"`, value, `"}`}, "")
	return []byte(ret)
}
func (o *OltpInf) startServer() error {
	o.echoSv = echo.New()
	o.echoSv.HideBanner = true
	o.echoSv.Use(ZapLogger(o.logger))
	o.echoSv.Use(middleware.Recover())

	// Routes
	o.echoSv.GET("/api/v1/status", o.getStatus)
	o.echoSv.GET("/api/v1/capabilities", o.getCapabilities)
	o.echoSv.GET("/api/v1/policies", o.getPolicies)
	o.echoSv.POST("/api/v1/policies", o.createPolicy)
	o.echoSv.GET("/api/v1/policies/:policy", o.getPolicy)
	o.echoSv.DELETE("/api/v1/policies/:policy", o.deletePolicy)

	serverHost := o.conf.OtlpInf.ServerHost
	serverPort := strconv.FormatUint(o.conf.OtlpInf.ServerPort, 10)
	o.echoSv.Logger.Fatal(o.echoSv.Start(serverHost + ":" + serverPort))
	return nil
}

func (o *OltpInf) getStatus(c echo.Context) error {
	o.stat.UpTime = time.Since(o.stat.StartTime)
	return c.JSONPretty(http.StatusOK, o.stat, "  ")
}

func (o *OltpInf) getCapabilities(c echo.Context) error {
	j, err := yson.YAMLToJSON(o.capabilities)
	if err != nil {
		return err
	}
	var ret interface{}
	err = json.Unmarshal(j, &ret)
	if err != nil {
		return err
	}
	return c.JSONPretty(http.StatusOK, ret, "  ")
}

func (o *OltpInf) getPolicies(c echo.Context) error {
	policies := make([]string, 0, len(o.policies))
	for k := range o.policies {
		policies = append(policies, k)
	}
	return c.JSONPretty(http.StatusOK, policies, "  ")
}

func (o *OltpInf) getPolicy(c echo.Context) error {
	policy := c.Param("policy")
	rInfo, ok := o.policies[policy]
	if ok {
		_, err := yaml.Marshal(rInfo.Policy)
		if err != nil {
			return c.JSONBlob(http.StatusBadRequest, messageReturn(err.Error()))
		}
		return c.Blob(http.StatusOK, "application/x-yaml", []byte(`à`))
	} else {
		return c.JSONBlob(http.StatusNotFound, messageReturn("Policy Not Found"))
	}
}

func (o *OltpInf) createPolicy(c echo.Context) error {
	if t := c.Request().Header.Get("Content-type"); t != "application/x-yaml" {
		return c.JSONBlob(http.StatusBadRequest,
			messageReturn("invalid Content-Type. Only 'application/x-yaml' is supported"))
	}
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	var payload map[string]config.Policy
	if err = yaml.Unmarshal(body, &payload); err != nil {
		return c.JSONBlob(http.StatusBadRequest, messageReturn(err.Error()))
	}
	if len(payload) > 1 {
		return c.JSONBlob(http.StatusBadRequest,
			messageReturn("only single policy allowed per request"))
	}
	var policy string
	var data config.Policy
	for policy, data = range payload {
		_, ok := o.policies[policy]
		if ok {
			return c.JSONBlob(http.StatusForbidden, messageReturn("policy already exists"))
		}
		if len(data.Config) == 0 {
			return c.JSONBlob(http.StatusForbidden, messageReturn("config field is required"))
		}
	}

	r := runner.New(o.logger, policy, o.policiesDir)
	if err := r.Configure(&data); err != nil {
		return c.JSONBlob(http.StatusBadRequest, messageReturn(err.Error()))
	}
	runnerCtx := context.WithValue(o.ctx, "routine", policy)
	if err := r.Start(context.WithCancel(runnerCtx)); err != nil {
		return c.JSONBlob(http.StatusBadRequest, messageReturn(err.Error()))
	}
	o.policies[policy] = RunnerInfo{Policy: data, Instance: r}
	return c.String(http.StatusOK, "Hello, World!")
}

func (o *OltpInf) deletePolicy(c echo.Context) error {
	policy := c.Param("policy")
	r, ok := o.policies[policy]
	if ok {
		if err := r.Instance.Stop(o.ctx); err != nil {
			return c.JSONBlob(http.StatusBadRequest, messageReturn(err.Error()))
		}
		delete(o.policies, policy)
		return c.JSONBlob(http.StatusOK, messageReturn(policy+" was deleted"))
	} else {
		return c.JSONBlob(http.StatusNotFound, messageReturn("Policy Not Found"))
	}
}
