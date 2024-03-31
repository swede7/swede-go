package example_test

import (
	"testing"
	"me.weldnor/swede/runner"
)

func addTwoAndTwo(ctx *runner.Context) runner.ScenarioExecutionResult {
	return runner.ScenarioExecutionResult{}
}

func checkResultIs4(ctx *runner.Context) runner.ScenarioExecutionResult {
	return runner.ScenarioExecutionResult{}
}
func TestSwedeRunner(t *testing.T) {
	_runner := runner.NewRunner()
	_runner.LoadFeatureFile("./feature/formatted.swede")
	_runner.RegisterFunc("", addTwoAndTwo)
	_runner.RegisterFunc("", checkResultIs4)
	_runner.RegisterFunc("", checkResultIs5)
	_runner.RegisterFunc("", checkResultIs6)
	_runner.Run()
}

func checkResultIs5(ctx *runner.Context) runner.ScenarioExecutionResult {
	return runner.ScenarioExecutionResult{}
}

func checkResultIs6(ctx *runner.Context) runner.ScenarioExecutionResult {
	return runner.ScenarioExecutionResult{}
}
//go:generate swede gen --feature-file ./feature/formatted.swede
// swede:step
// swede:step
// swede:step
// swede:step
