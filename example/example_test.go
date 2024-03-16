package example_test

import (
	"me.weldnor/swede/runner"
	"testing"
)

//go:generate swede gen --feature-file=./feature/formatted.swede

// swede:step Add 2 and 2
func addTwoAndTwo(ctx *runner.ScenarioContext) runner.ScenarioExecutionResult {
	ctx.Variables["result"] = 2 + 2

	return runner.ScenarioExecutionResult{Status: true}
}

func checkResultIs4(ctx *runner.ScenarioContext) runner.ScenarioExecutionResult {
	result := ctx.Variables["result"].(int)

	if result != 4 {
		return runner.ScenarioExecutionResult{false, "result is not 4"}
	}

	return runner.ScenarioExecutionResult{Status: true}
}

// swede:step Check that result is 5
func checkResultIs5(ctx *runner.ScenarioContext) runner.ScenarioExecutionResult {
	result := ctx.Variables["result"].(int)

	if result != 5 {
		return runner.ScenarioExecutionResult{false, "result is not 5"}
	}

	return runner.ScenarioExecutionResult{Status: true}
}

func checkResultIs6(ctx *runner.ScenarioContext) runner.ScenarioExecutionResult {
	result := ctx.Variables["result"].(int)

	if result != 5 {
		return runner.ScenarioExecutionResult{false, "result is not 5"}
	}

	return runner.ScenarioExecutionResult{Status: true}
}
func TestSwedeRunner(t *testing.T) {
	_runner :=
		runner.
			NewRunner()
	_runner.LoadFeatureFile("./feature/formatted.swede")

	_runner.RegisterFunc("Add 2 and 2", addTwoAndTwo)
	_runner.RegisterFunc("Check that result is 4",
		checkResultIs4)
	_runner.RegisterFunc("Check that result is 5",
		checkResultIs5)
	_runner.Run()
}
