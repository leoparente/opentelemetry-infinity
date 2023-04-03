package runner

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"regexp"
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

type RunningStatus int

const (
	Unknown RunningStatus = iota
	Running
	RunnerError
	Offline
)

var MapRunningStatus = map[RunningStatus]string{
	Unknown:     "unknown",
	Running:     "running",
	RunnerError: "runner_error",
	Offline:     "offline",
}

type State struct {
	Status            RunningStatus `yaml:"status"`
	startTime         time.Time
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
	state        State
	cancelFunc   context.CancelFunc
	ctx          context.Context
	cmd          *exec.Cmd
	errChan      chan string
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
	return Runner{logger: logger, policyName: policyName, policyDir: policyDir, errChan: make(chan string)}
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
	r.cmd.Stderr = &RunnerStderr{r.logger, r.policyName, r.errChan}
	if err = r.cmd.Start(); err != nil {
		return err
	}
	reg, err := regexp.Compile("[^a-zA-Z0-9:(), ]+")
	if err != nil {
		return err
	}
	r.state.startTime = time.Now()
	ctxTimeout, cancel := context.WithTimeout(r.ctx, 1*time.Second)
	defer cancel()
	select {
	case line := <-r.errChan:
		return errors.New(string(append([]byte("otelcol-contrib - "), reg.ReplaceAllString(line, "")...)))
	case <-ctxTimeout.Done():
		r.state.Status = Running
		r.logger.Info("runner proccess started successfully", zap.String("policy", r.policyName), zap.Any("pid", r.cmd.Process.Pid))
	}

	go func() {
		for {
			select {
			case line := <-r.errChan:
				r.state.LastError = string(append([]byte("otelcol-contrib - "), reg.ReplaceAllString(line, "")...))
				r.state.Status = RunnerError
			case <-r.ctx.Done():
				r.Stop(r.ctx)
			}
		}
	}()

	return nil
}

func (r *Runner) Stop(ctx context.Context) error {
	r.logger.Info("routine call to stop runner", zap.Any("routine", ctx.Value("routine")))
	defer r.cancelFunc()
	if err := r.cmd.Cancel(); err != nil {
		return err
	}
	r.state.Status = Offline
	pid := r.cmd.ProcessState.Pid()
	exitCode := r.cmd.ProcessState.ExitCode()
	r.logger.Info("runner process stopped", zap.Int("pid", pid), zap.Int("exit_code", exitCode))
	return nil
}

func (r *Runner) FullReset(ctx context.Context) error {
	return nil
}

func (r *Runner) GetRunningStatus() (RunningStatus, string, error) {
	return r.state.Status, r.state.LastError, nil
}
