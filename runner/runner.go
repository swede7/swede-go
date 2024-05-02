package runner

import (
	"errors"
	"fmt"
	"strings"

	"me.weldnor/swede/core/lang/swede/parser"
)

type Runner struct {
	registeredFuncs map[string]StepFunc

	beforeScenarioFunc HandlerFunc
	afterScenarioFunc  HandlerFunc

	beforeFeatureFunc HandlerFunc
	afterFeatureFunc  HandlerFunc

	featureName string
	scenarios   []Scenario
}

type StepFunc func(*Context) error
type HandlerFunc func(*Context)

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

func (r *Runner) RegisterFunc(name string, registeredFunc StepFunc) {
	if r.registeredFuncs == nil {
		r.registeredFuncs = make(map[string]StepFunc)
	}

	r.registeredFuncs[name] = registeredFunc
}

func (r *Runner) RegisterBeforeScenarioFunc(registeredFunc HandlerFunc) {
	r.beforeScenarioFunc = registeredFunc
}

func (r *Runner) RegisterAfterScenarioFunc(registeredFunc HandlerFunc) {
	r.afterScenarioFunc = registeredFunc
}

func (r *Runner) RegisterBeforeFeatureFunc(registeredFunc HandlerFunc) {
	r.beforeFeatureFunc = registeredFunc
}

func (r *Runner) RegisterAfterFeatureFunc(registeredFunc HandlerFunc) {
	r.afterFeatureFunc = registeredFunc
}

func (r *Runner) Run() {
	fmt.Printf("Running feature %s\n\n", r.featureName)

	context := newContext()

	if r.beforeFeatureFunc != nil {
		r.beforeFeatureFunc(context)
	}

	for _, scenario := range r.scenarios {
		r.executeScenario(scenario, context)
	}

	if r.afterFeatureFunc != nil {
		r.afterFeatureFunc(context)
	}
}

type executionStatus string

const (
	passed  executionStatus = "passed"
	failed  executionStatus = "failed"
	skipped executionStatus = "skipped"
)

func (r *Runner) executeScenario(scenario Scenario, context *Context) {
	fmt.Printf("\tRunning scenario %s\n\n", scenario.Name)

	if r.beforeScenarioFunc != nil {
		r.beforeScenarioFunc(context)
	}

	isFailed := false

	for _, step := range scenario.Steps {
		if isFailed {
			fmt.Printf("\t\t- Skipping step %s ⚠️\n", step.Name)
			continue
		}

		fmt.Printf("\t\t- Running step %s ", step.Name)

		result := r.executeStep(step, context)

		switch result.status {
		case passed:
			fmt.Println("✅")
		case failed:
			fmt.Println("❌")
			fmt.Printf(" failed with error: %s\n", result.error)
			isFailed = true
		}

	}

	if r.afterScenarioFunc != nil {
		r.afterScenarioFunc(context)
	}

	fmt.Println()
}

type stepExecutionResult struct {
	status executionStatus
	error  error
}

func (r *Runner) executeStep(step Step, context *Context) stepExecutionResult {
	stepName := step.Name
	registeredFunc, ok := r.registeredFuncs[stepName]

	if !ok {
		return stepExecutionResult{status: failed, error: errors.New("step is not defined")}
	}

	err := registeredFunc(context)

	if err != nil {
		return stepExecutionResult{status: failed, error: err}
	}

	return stepExecutionResult{status: passed}
}
