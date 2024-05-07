package example_test

import (
	"errors"
	"fmt"
	"me.weldnor/swede/runner"
	"testing"
	"time"
)

//go:generate swede gen -f ./feature/calculator.feature

// swede:step Add <first:int> and <second:int>
func addIntAndInt(ctx *runner.Context) error {
	result := ctx.GetIntVariable("first") + ctx.GetIntVariable("second")
	ctx.SetVariable("result", result)
	return nil
}

// swede:step Add string <first:string> and <second:string>
func addStringAndString(ctx *runner.Context) error {
	result := ctx.GetStringVariable("first") + ctx.GetStringVariable("second")
	ctx.SetVariable("result", result)
	return nil
}

// swede:step Check that result is <expected:int>
func checkResultIs(ctx *runner.Context) error {
	result := ctx.GetIntVariable("result")

	expected := ctx.GetIntVariable("expected")

	if result != expected {
		return errors.New(fmt.Sprintf("result != %d", expected))
	}

	return nil
}

// swede:step Check that result string is <expected:string>
func checkResultStringIs(ctx *runner.Context) error {
	result := ctx.GetStringVariable("result")

	expected := ctx.GetStringVariable("expected")

	if result != expected {
		return errors.New(fmt.Sprintf("result != %s", expected))
	}

	return nil
}

// swede:beforeScenario
func setUp(ctx *runner.Context) {
	fmt.Println("setting up environment...")
	time.Sleep(1 * time.Second)
}

// swede:afterScenario
func tearDown(ctx *runner.Context) {
	fmt.Println("clearing up...")
	time.Sleep(1 * time.Second)
}

func TestSwedeRunner(t *testing.T) {
	_runner := runner.NewRunner()
	_runner.LoadFeatureFile("./feature/calculator.feature")
	_runner.RegisterFunc("Add <first:int> and <second:int>", addIntAndInt)
	_runner.RegisterFunc("Add string <first:string> and <second:string>", addStringAndString)
	_runner.RegisterFunc("Check that result is <expected:int>", checkResultIs)
	_runner.RegisterFunc("Check that result string is <expected:string>", checkResultStringIs)
	_runner.RegisterBeforeScenarioFunc(setUp)
	_runner.RegisterAfterScenarioFunc(tearDown)
	_runner.Run()
}
