package otlpinf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ghodss/yaml"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leoparente/otlpinf/config"
	"github.com/leoparente/otlpinf/runner"
)

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
	j, err := yaml.YAMLToJSON(o.capabilities)
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
	_, ok := o.policies[policy]
	if ok {
		return c.Blob(http.StatusOK, "application/x-yaml", []byte(`à`))
	} else {
		return c.JSONBlob(http.StatusNotFound, []byte(`{"message":"Policy Not Found"}`))
	}
}

func (o *OltpInf) createPolicy(c echo.Context) error {
	if t := c.Request().Header.Get("Content-type"); t != "application/x-yaml" {
		return c.JSONBlob(http.StatusBadRequest, []byte(`{"message": 
			"invalid Content-Type. Only 'application/x-yaml' is supported"}`))
	}
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	var payload map[string]config.Policy
	if err = yaml.y(body, &payload); err != nil {
		o.logger.Info(fmt.Sprint("%v", err))
		return err
	}

	o.logger.Info(fmt.Sprint("%v", payload))

	r := runner.New(o.logger, "test")
	runnerCtx := context.WithValue(o.ctx, "routine", "test")
	if err := r.Start(context.WithCancel(runnerCtx)); err != nil {
		return err
	}
	// o.runners["test"] = r
	// o.runnerState["test"] = &runner.State{
	// 	Status:        runner.Unknown,
	// 	LastRestartTS: time.Now(),
	// }
	return c.String(http.StatusOK, "Hello, World!")
}

func (o *OltpInf) deletePolicy(c echo.Context) error {
	policy := c.Param("policy")
	_, ok := o.policies[policy]
	if ok {
		return c.String(http.StatusOK, "Hello, World!")
	} else {
		return c.JSONBlob(http.StatusNotFound, []byte(`{"message":"Policy Not Found"}`))
	}
}
