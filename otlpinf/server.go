package otlpinf

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	o.echoSv.DELETE("/api/v1/policies/:name", o.getPolicy)
	o.echoSv.DELETE("/api/v1/policies/:name", o.deletePolicy)

	serverHost := o.conf.OtlpInf.ServerHost
	serverPort := strconv.FormatUint(o.conf.OtlpInf.ServerPort, 10)
	o.echoSv.Logger.Fatal(o.echoSv.Start(serverHost + ":" + serverPort))
	return nil
}

func (o *OltpInf) getStatus(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (o *OltpInf) getCapabilities(c echo.Context) error {
	return c.Blob(http.StatusOK, "text/x-yaml", o.capabilities)
}

func (o *OltpInf) getPolicies(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (o *OltpInf) getPolicy(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (o *OltpInf) createPolicy(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (o *OltpInf) deletePolicy(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
