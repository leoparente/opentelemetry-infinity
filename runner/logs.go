package runner

import "go.uber.org/zap"

type RunnerStdout struct {
	logger     *zap.Logger
	policyName string
}

type RunnerStderr struct {
	logger     *zap.Logger
	policyName string
}

func (rs *RunnerStdout) Write(p []byte) (n int, err error) {
	rs.logger.Info("otelcol-contrib stdout", zap.String("policy", rs.policyName), zap.ByteString("log", p))
	n = len(p)
	return
}

func (rs *RunnerStderr) Write(p []byte) (n int, err error) {
	rs.logger.Error("runner stderr", zap.String("policy", rs.policyName), zap.ByteString("log", p))
	n = len(p)
	return
}
