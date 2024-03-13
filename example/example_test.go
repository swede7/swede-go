package example_test

import (
	"testing"

	"me.weldnor/swede/runner"
)

//go:generate swede gen

// swede:step
func anotherFunc() {
	// do nothing
}

func TestLexerForCodeExample(t *testing.T) {
	_runner := runner.NewRunner()
	_runner.LoadFeatureFile("./feature/formatted.swede")

	_runner.RegisterFunc("Add \"2\" and \"2\"", func(ctx *runner.ScenarioContext) runner.ScenarioExecutionResult {
		ctx.Variables["sum"] = 2 + 2
		return runner.ScenarioExecutionResult{Status: true, Message: "OK"}
	})

	_runner.RegisterFunc("Check that result is \"4\"", func(ctx *runner.ScenarioContext) runner.ScenarioExecutionResult {
		sum := ctx.Variables["sum"].(int)
		if sum != 4 {
			return runner.ScenarioExecutionResult{Status: false, Message: "oops! sum is not equal to 4"}
		}

		return runner.ScenarioExecutionResult{Status: true, Message: "OK"}
	})

	_runner.RegisterFunc("Check that result is \"5\"", func(ctx *runner.ScenarioContext) runner.ScenarioExecutionResult {
		sum := ctx.Variables["sum"].(int)
		if sum != 5 {
			return runner.ScenarioExecutionResult{Status: false, Message: "oops! sum is not equal to 5"}
		}

		return runner.ScenarioExecutionResult{Status: true, Message: "OK"}
	})

	_runner.Run()
}
