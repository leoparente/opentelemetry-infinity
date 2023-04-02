package otlpinf

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/leoparente/otlpinf/config"
	"github.com/leoparente/otlpinf/runner"
	"go.uber.org/zap"
)

type OltpInf struct {
	logger         *zap.Logger
	conf           *config.Config
	policies       []string
	runners        map[string]runner.Runner
	runnerState    map[string]*runner.State
	cancelFunction context.CancelFunc
	echoSv         *echo.Echo
	capabilities   []byte
}

func New(logger *zap.Logger, c *config.Config) (OltpInf, error) {
	return OltpInf{logger: logger, conf: c}, nil
}

func (o *OltpInf) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	otlpInfCtx := context.WithValue(ctx, "routine", "otlpInfRoutine")
	o.cancelFunction = cancelFunc
	var err error
	o.capabilities, err = runner.GetCapabilities()
	if err != nil {
		return err
	}
	if err = o.startServer(); err != nil {
		return err
	}
	r := runner.New(o.logger, "test")
	runnerCtx := context.WithValue(otlpInfCtx, "routine", "test")
	if err := r.Start(context.WithCancel(runnerCtx)); err != nil {
		return err
	}
	// o.runners["test"] = r
	// o.runnerState["test"] = &runner.State{
	// 	Status:        runner.Unknown,
	// 	LastRestartTS: time.Now(),
	// }

	return nil
}

func (o *OltpInf) Stop(ctx context.Context) {
	o.logger.Info("routine call for stop agent", zap.Any("routine", ctx.Value("routine")))
	for name, b := range o.runners {
		if state, _, _ := b.GetRunningStatus(); state == runner.Running {
			o.logger.Debug("stopping runner", zap.String("runner", name))
			if err := b.Stop(ctx); err != nil {
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
