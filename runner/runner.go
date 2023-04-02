package runner

import (
	"context"
	"os"

	_ "embed"
	"time"

	"github.com/amenzhinsky/go-memexec"
	"go.uber.org/zap"
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
	logger     *zap.Logger
	policyName string
	startTime  time.Time
	cancelFunc context.CancelFunc
	ctx        context.Context
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

func New(logger *zap.Logger, policyName string) Runner {
	return Runner{logger: logger, policyName: policyName}
}

func (r *Runner) getProcRunningStatus() (RunningStatus, string, error) {
	// status := r.proc.Status()
	// if status.Error != nil {
	// 	errMsg := fmt.Sprintf("runner process error: %v", status.Error)
	// 	return RunnerError, errMsg, status.Error
	// }
	// if status.Complete {
	// 	err := r.proc.Stop()
	// 	return Offline, "runner process ended", err
	// }
	// if status.StopTs > 0 {
	// 	return Offline, "runner process ended", nil
	// }
	return Running, "", nil
}

func (r *Runner) Version() (string, error) {
	return "", nil
}

func (r *Runner) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	r.startTime = time.Now()
	r.cancelFunc = cancelFunc
	r.ctx = ctx

	f, err := os.CreateTemp("", r.policyName)
	if err != nil {
		return err
	}
	f.Write([]byte("ahehaehae"))
	f.Close()
	r.logger.Info(f.Name())

	sOptions := []string{
		"--config",
		f.Name(),
	}
	// fil, err := os.Open("src/a/b")
	// if err != nil {
	// 	r.logger.Error("error", zap.Error(err))
	// 	return err
	// }
	// r.logger.Error("error" + fil.Name())
	exe, err := memexec.New(otel_contrib)
	if err != nil {
		return err
	}
	defer exe.Close()

	sample := &RunnerStdout{r.logger, r.policyName}
	cmd := exe.CommandContext(ctx, sOptions...)
	cmd.Stdout = sample
	cmd.Stderr = &RunnerStderr{r.logger, r.policyName}
	cmd.Start()

	cmd.Wait()
	// data, err := cmd.Output() // cmd is a `*exec.Cmd` from the standard libraryp
	// if err != nil {
	// 	r.logger.Info("erro", zap.Error(err))
	// }
	// r.logger.Info(string(data))

	return nil
}

func (r *Runner) Stop(ctx context.Context) error {
	// r.logger.Info("routine call to stop runner", zap.Any("routine", ctx.Value("routine")))
	// defer r.cancelFunc()
	// err := r.proc.Stop()
	// finalStatus := <-r.statusChan
	// if err != nil {
	// 	r.logger.Error("runner shutdown error", zap.Error(err))
	// 	return err
	// }
	// r.logger.Info("runner process stopped", zap.Int("pid", finalStatus.PID), zap.Int("exit_code", finalStatus.Exit))
	return nil
}

func (r *Runner) FullReset(ctx context.Context) error {
	return nil
}

func (r *Runner) GetStartTime() time.Time {
	return r.startTime
}

func (r *Runner) GetRunningStatus() (RunningStatus, string, error) {
	runningStatus, errMsg, err := r.getProcRunningStatus()
	if runningStatus != Running {
		return runningStatus, errMsg, err
	}
	return runningStatus, "", nil
}
