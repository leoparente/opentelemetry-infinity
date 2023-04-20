package otlpinf

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoparente/otlpinf/config"
	"github.com/leoparente/otlpinf/runner"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type RunnerInfo struct {
	Policy   config.Policy
	Instance runner.Runner
}

type OltpInf struct {
	logger         *zap.Logger
	conf           *config.Config
	stat           config.Status
	policies       map[string]RunnerInfo
	policiesDir    string
	ctx            context.Context
	cancelFunction context.CancelFunc
	router         *gin.Engine
	capabilities   []byte
}

func New(logger *zap.Logger, c *config.Config) (OltpInf, error) {
	return OltpInf{logger: logger, conf: c, policies: make(map[string]RunnerInfo)}, nil
}

func (o *OltpInf) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	o.stat.StartTime = time.Now()
	o.stat.InfVersion = o.conf.Version
	o.ctx = context.WithValue(ctx, "routine", "otlpInfRoutine")
	o.cancelFunction = cancelFunc

	var err error
	o.policiesDir, err = os.MkdirTemp("", "policies")
	if err != nil {
		return err
	}
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
	defer os.RemoveAll(o.policiesDir)
	o.cancelFunction()
}

func (o *OltpInf) RestartRunner(ctx context.Context, name string, reason string) error {
	return nil
}

func (o *OltpInf) RestartAll(ctx context.Context, reason string) error {
	return nil
}
