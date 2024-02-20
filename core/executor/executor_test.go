package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutor_RegisterFunc(t *testing.T) {
	e := Executor{}
	e.RegisterFunc("testFunc", func(scenarioContext *ScenarioContext) StepExecutionResult {
		return StepExecutionResult{Status: true, Message: "test"}
	})

	assert.NotNil(t, e.registeredFuncs["testFunc"], "Expected function to be registered, but it was not")
}

func TestExecutor_ExecuteScenario(t *testing.T) {
	e := Executor{}
	e.RegisterFunc("testFunc", func(scenarioContext *ScenarioContext) StepExecutionResult {
		return StepExecutionResult{Status: true, Message: "test"}
	})

	scenario := Scenario{
		Name: "testScenario",
		Steps: []Step{
			{Name: "testFunc"},
		},
	}

	result := e.ExecuteScenario(scenario)
	assert.True(t, result.Status, "Expected scenario to execute successfully, but got: %v", result)
	assert.Equal(t, "ok", result.Message, "Expected scenario to execute successfully, but got: %v", result)
}

func TestExecutor_ExecuteScenariosParallel(t *testing.T) {
	e := Executor{}
	e.RegisterFunc("testFunc", func(scenarioContext *ScenarioContext) StepExecutionResult {
		return StepExecutionResult{Status: true, Message: "test"}
	})

	scenarios := []Scenario{
		{
			Name: "testScenario1",
			Steps: []Step{
				{Name: "testFunc"},
			},
		},
		{
			Name: "testScenario2",
			Steps: []Step{
				{Name: "testFunc"},
			},
		},
	}

	e.scenarios = scenarios
	e.ExecuteScenariosParallel()

	// Since the ExecuteScenariosParallel method prints to the console, we can't directly
	// check the results. However, we can check if the method completes without panics.
}
