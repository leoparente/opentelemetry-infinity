package runner

import (
	"bytes"
	"io"

	"go.uber.org/zap"
)

type RunnerOutput struct {
	err        bool
	logger     *zap.Logger
	policyName string
	channel    chan<- string
}

var _ io.Writer = (*RunnerOutput)(nil)

func (rs *RunnerOutput) Write(p []byte) (n int, err error) {
	if rs.err {
		rs.logger.Error("otelcol-contrib stderr", zap.String("policy", rs.policyName), zap.String("log", bytes.NewBuffer(p).String()))
		rs.channel <- string(p)
	} else {
		rs.logger.Info("otelcol-contrib stdout", zap.String("policy", rs.policyName), zap.Any("log", bytes.NewBuffer(p).String()))
	}
	return len(p), nil
}
