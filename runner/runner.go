package runner

import (
	"context"
	"os"
	"os/exec"
	"strings"

	_ "embed"
	"time"

	"github.com/amenzhinsky/go-memexec"
	"github.com/leoparente/otlpinf/config"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

//go:embed otelcol-contrib
var otel_contrib []byte

const (
	Unknown RunningStatus = iota
	Running
	RunnerError
	Offline
)

type RunningStatus int

type State struct {
	Status            RunningStatus
	RestartCount      int64
	LastError         string
	LastRestartTS     time.Time
	LastRestartReason string
}

type Runner struct {
	logger       *zap.Logger
	policyName   string
	policyDir    string
	policyFile   string
	featureGates string
	startTime    time.Time
	cancelFunc   context.CancelFunc
	ctx          context.Context
	cmd          *exec.Cmd
}

func GetCapabilities() ([]byte, error) {
	exe, err := memexec.New(otel_contrib)
	if err != nil {
		return nil, err
	}
	defer exe.Close()
	cmd := exe.Command("components")
	ret, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func New(logger *zap.Logger, policyName string, policyDir string) Runner {
	return Runner{logger: logger, policyName: policyName, policyDir: policyDir}
}

func (r *Runner) Configure(c *config.Policy) error {
	b, err := yaml.Marshal(&c.Config)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(r.policyDir, r.policyName)
	if err != nil {
		return err
	}
	if _, err = f.Write(b); err != nil {
		return err
	}
	r.policyFile = f.Name()
	if err = f.Close(); err != nil {
		return err
	}

	if c.FeatureGates != nil {
		r.featureGates = strings.Join(c.FeatureGates, ",")
	}

	return nil
}

func (r *Runner) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	r.startTime = time.Now()
	r.cancelFunc = cancelFunc
	r.ctx = ctx

	sOptions := []string{
		"--config",
		r.policyFile,
	}

	if len(r.featureGates) > 0 {
		sOptions = append(sOptions, "--feature-gates", r.featureGates)
	}

	//TODO: implement set support

	exe, err := memexec.New(otel_contrib)
	if err != nil {
		return err
	}
	defer exe.Close()

	r.cmd = exe.CommandContext(ctx, sOptions...)
	r.cmd.Stdout = &RunnerStdout{r.logger, r.policyName}
	r.cmd.Stderr = &RunnerStderr{r.logger, r.policyName}
	if err = r.cmd.Start(); err != nil {
		return err
	}

	// data, err := cmd.Output() // cmd is a `*exec.Cmd` from the standard libraryp
	// if err != nil {
	// 	r.logger.Info("erro", zap.Error(err))
	// }
	// r.logger.Info(string(data))

	return nil
}

func (r *Runner) Stop(ctx context.Context) error {
	r.logger.Info("routine call to stop runner", zap.Any("routine", ctx.Value("routine")))
	defer r.cancelFunc()
	if err := r.cmd.Cancel(); err != nil {
		return err
	}
	pid := r.cmd.ProcessState.Pid()
	exitCode := r.cmd.ProcessState.ExitCode()
	r.logger.Info("runner process stopped", zap.Int("pid", pid), zap.Int("exit_code", exitCode))
	return nil
}

func (r *Runner) FullReset(ctx context.Context) error {
	return nil
}

func (r *Runner) GetRunningStatus() (RunningStatus, string, error) {
	// runningStatus, errMsg, err := r.getProcRunningStatus()
	// if runningStatus != Running {
	// 	return runningStatus, errMsg, err
	// }
	return 2, "", nil
}
