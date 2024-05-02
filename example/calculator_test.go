package example_test

import (
	"errors"
	"fmt"
	"me.weldnor/swede/runner"
	"testing"
	"time"
)

//go:generate swede gen -f ./feature/calculator.feature

// swede:step Add 2 and 2
func addTwoAndTwo(ctx *runner.Context) error {
	result := 2 + 2
	ctx.SetVariable("result", result)

	return nil
}

// swede:step Check that result is 4
func checkResultIs4(ctx *runner.Context) error {
	result := ctx.GetVariable("result").(int)

	if result != 4 {
		return errors.New("result != 4")
	}

	return nil
}

// swede:beforeScenario
func setUp(ctx *runner.Context) {
	fmt.Println("setting up environment...")
	time.Sleep(3 * time.Second)
}

// swede:afterScenario
func tearDown(ctx *runner.Context) {
	fmt.Println("clearing up...")
	time.Sleep(3 * time.Second)
}

func TestSwedeRunner(t *testing.T) {
	_runner := runner.NewRunner()
	_runner.LoadFeatureFile("./feature/calculator.feature")
	_runner.RegisterFunc("Add 2 and 2", addTwoAndTwo)
	_runner.RegisterFunc("Check that result is 4", checkResultIs4)
	_runner.RegisterBeforeScenarioFunc(setUp)
	_runner.RegisterAfterScenarioFunc(tearDown)
	_runner.Run()
}
