package otlpinf

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/leoparente/otlpinf/config"
	"github.com/leoparente/otlpinf/runner"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type RunnerInfo struct {
	Policy   config.Policy
	Instance runner.Runner
	State    runner.State
}

type OltpInf struct {
	logger         *zap.Logger
	conf           *config.Config
	stat           config.Status
	policies       map[string]RunnerInfo
	ctx            context.Context
	cancelFunction context.CancelFunc
	echoSv         *echo.Echo
	capabilities   []byte
}

func New(logger *zap.Logger, c *config.Config) (OltpInf, error) {
	return OltpInf{logger: logger, conf: c}, nil
}

func (o *OltpInf) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	o.stat.StartTime = time.Now()
	o.stat.InfVersion = o.conf.Version
	o.ctx = context.WithValue(ctx, "routine", "otlpInfRoutine")
	o.cancelFunction = cancelFunc
	var err error
	o.capabilities, err = runner.GetCapabilities()
	if err != nil {
		return err
	}
	s := struct {
		Buildinfo struct {
			Version string
		}
	}{}
	err = yaml.Unmarshal(o.capabilities, &s)
	if err != nil {
		return err
	}
	o.stat.ContribVersion = s.Buildinfo.Version
	if err = o.startServer(); err != nil {
		return err
	}
	return nil
}

func (o *OltpInf) Stop(ctx context.Context) {
	o.logger.Info("routine call for stop otlpinf", zap.Any("routine", ctx.Value("routine")))
	for name, b := range o.policies {
		if state, _, _ := b.Instance.GetRunningStatus(); state == runner.Running {
			o.logger.Debug("stopping runner", zap.String("runner", name))
			if err := b.Instance.Stop(ctx); err != nil {
				o.logger.Error("error while stopping the runner", zap.String("runner", name))
			}
		}
	}
	defer o.cancelFunction()
}

func (o *OltpInf) RestartRunner(ctx context.Context, name string, reason string) error {
	return nil
}

func (o *OltpInf) RestartAll(ctx context.Context, reason string) error {
	return nil
}
