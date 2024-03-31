package runner

import (
	"fmt"
	"strings"

	"me.weldnor/swede/core/lang/swede/parser"
)

type Runner struct {
	registeredFuncs map[string]RegisteredFunc

	beforeScenarioFunc RegisteredFunc
	afterScenarioFunc  RegisteredFunc

	beforeFeatureFunc RegisteredFunc
	afterFeatureFunc  RegisteredFunc

	featureName string
	scenarios   []Scenario
}

type RegisteredFunc func(*Context) ScenarioExecutionResult

func NewRunner() *Runner {
	return new(Runner)
}

func (r *Runner) LoadFeatureFile(path string) {
	parserResult := parser.ParseFile(path)

	if parserResult.Errors != nil {
		panic("feature file contains errors")
	}

	rootNode := parserResult.RootNode
	r.featureName = strings.TrimSpace(rootNode.GetChildByType(parser.FEATURE).Value)

	scenarios := make([]Scenario, 0)

	for _, scenarioNode := range rootNode.GetChildrenByType(parser.SCENARIO) {
		scenario := Scenario{}
		scenario.Name = strings.TrimSpace(scenarioNode.Value)

		for _, stepNode := range scenarioNode.GetChildrenByType(parser.STEP) {
			step := Step{}
			step.Name = strings.TrimSpace(stepNode.Value)
			scenario.Steps = append(scenario.Steps, step)
		}

		scenarios = append(scenarios, scenario)
	}

	r.scenarios = scenarios
}

func (e *Runner) RegisterFunc(name string, registeredFunc RegisteredFunc) {
	if e.registeredFuncs == nil {
		e.registeredFuncs = make(map[string]RegisteredFunc)
	}

	e.registeredFuncs[name] = registeredFunc
}

func (e *Runner) RegisterBeforeScenarioFunc(registeredFunc RegisteredFunc) {
	e.beforeScenarioFunc = registeredFunc
}

func (e *Runner) RegisterAfterScenarioFunc(registeredFunc RegisteredFunc) {
	e.afterScenarioFunc = registeredFunc
}

func (e *Runner) RegisterBeforeFeatureFunc(registeredFunc RegisteredFunc) {
	e.beforeFeatureFunc = registeredFunc
}

func (e *Runner) RegisterAfterFeatureFunc(registeredFunc RegisteredFunc) {
	e.afterFeatureFunc = registeredFunc
}

func (e *Runner) Run() {
	fmt.Printf("start running %s\n", e.featureName)

	context := newContext()

	if e.beforeFeatureFunc != nil {
		e.beforeFeatureFunc(context)
	}

	for _, scenario := range e.scenarios {
		result := e.executeScenario(scenario, context)

		if result.Status {
			fmt.Printf("[PASS] %s\n", scenario.Name)
		} else {
			fmt.Printf("[FAIL] %s: %s\n", scenario.Name, result.Message)
		}
	}

	if e.afterFeatureFunc != nil {
		e.afterFeatureFunc(context)
	}
}

func (e *Runner) executeScenario(scenario Scenario, context *Context) ScenarioExecutionResult {
	if e.beforeScenarioFunc != nil {
		e.beforeScenarioFunc(context)
	}

	for _, step := range scenario.Steps {
		stepName := step.Name
		registeredFunc, ok := e.registeredFuncs[stepName]

		if !ok {
			return ScenarioExecutionResult{false, "step with name " + stepName + " is not registered"}
		}

		stepResult := registeredFunc(context)

		if !stepResult.Status {
			return ScenarioExecutionResult{false, "failed on step " + stepName}
		}
	}

	if e.afterScenarioFunc != nil {
		e.afterScenarioFunc(context)
	}

	return ScenarioExecutionResult{Status: true}
}
