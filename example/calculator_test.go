package example_test

import (
	"me.weldnor/swede/runner"
	"testing"
)

//go:generate swede gen -f ./feature/calculator.feature

// swede:step add 2 and 2
func addTwoAndTwo(ctx *runner.Context) runner.ScenarioExecutionResult {
	result := 2 + 2
	ctx.SetVariable("result", result)

	return runner.ScenarioExecutionResult{
		Status: true,
	}
}

// swede:step check result is 4
func checkResultIs4(ctx *runner.Context) runner.ScenarioExecutionResult {
	result := ctx.GetVariable("result").(int)

	if result != 4 {
		return runner.ScenarioExecutionResult{
			Status:  false,
			Message: "not equals",
		}
	}

	return runner.ScenarioExecutionResult{
		Status: true,
	}
}

func TestSwedeRunner(t *testing.T) {
	_runner := runner.NewRunner()
	_runner.LoadFeatureFile("./feature/calculator.feature")
	_runner.RegisterFunc("add 2 and 2", addTwoAndTwo)
	_runner.RegisterFunc("check result is 4", checkResultIs4)

	_runner.Run()
}
