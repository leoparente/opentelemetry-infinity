package runner

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/leoparente/opentelemetry-infinity/config"
	"go.uber.org/zap/zaptest"
	"gopkg.in/yaml.v2"
)

func TestRunnerNew(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	policyName := "test-policy"
	policyDir := "/tmp"
	selfTelemetry := false

	// Act
	runner := New(logger, policyName, policyDir, selfTelemetry)

	// Assert
	if runner.logger != logger {
		t.Errorf("Expected logger to be set to %v, got %v", logger, runner.logger)
	}

	if runner.policyName != policyName {
		t.Errorf("Expected policyName to be set to %s, got %s", policyName, runner.policyName)
	}

	if runner.policyDir != policyDir {
		t.Errorf("Expected policyDir to be set to %s, got %s", policyDir, runner.policyDir)
	}

	if runner.selfTelemetry != selfTelemetry {
		t.Errorf("Expected selfTelemetry to be set to %v, got %v", selfTelemetry, runner.selfTelemetry)
	}

	if len(runner.sets) != 0 {
		t.Errorf("Expected sets to be an empty slice, got %v", runner.sets)
	}
}

func TestRunnerConfigure(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	policyName := "test-policy"
	policyDir := "/tmp"
	enableTelemetry := true
	runner := &Runner{
		logger:        logger,
		policyName:    policyName,
		policyDir:     policyDir,
		selfTelemetry: enableTelemetry,
	}
	config := &config.Policy{
		FeatureGates: []string{"gate1", "gate2"},
		Set: map[string]string{
			"set1": "set1",
		},
		Config: map[string]interface{}{
			"policy": "value1",
		},
	}

	// Act
	err := runner.Configure(config)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	expectedFeatureGates := "gate1,gate2"
	if !reflect.DeepEqual(runner.featureGates, expectedFeatureGates) {
		t.Errorf("Expected featureGates to be %v, but got %v", expectedFeatureGates, runner.featureGates)
	}

	expectedSet := []string{"--set=set1=set1"}
	if !reflect.DeepEqual(runner.sets, expectedSet) {
		t.Errorf("Expected set to be %v, but got %v", expectedSet, runner.sets)
	}

	if !strings.Contains(runner.policyFile, policyName) {
		t.Errorf("Expected policy File to contain %v, but got %v", policyName, runner.policyFile)
	}
}

func TestRunnerStartStop(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	policyName := "test-policy"
	policyDir := "/tmp"
	enableTelemetry := true
	runner := &Runner{
		logger:        logger,
		policyName:    policyName,
		policyDir:     policyDir,
		selfTelemetry: enableTelemetry,
	}
	config := &config.Policy{
		FeatureGates: []string{"gate1", "gate2"},
		Set: map[string]string{
			"set1": "set1",
			"set2": "set2",
		},
		Config: map[string]interface{}{
			"policy": "value1",
		},
	}

	//Act
	err := runner.Configure(config)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	err = runner.Start(ctx, cancel)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	runner.Stop(ctx)

	s := runner.GetStatus()
	if MapStatus[s.Status] != "offline" {
		t.Errorf("Expected status to be offline, but got %v", MapStatus[s.Status])
	}
}

func TestRunnerGetCapabilities(t *testing.T) {
	//Act
	caps, err := GetCapabilities()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Assert
	s := struct {
		Buildinfo struct {
			Version string
		}
	}{}
	err = yaml.Unmarshal(caps, &s)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}
