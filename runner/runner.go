package runner

import (
	"fmt"
	"strings"
	"sync"

	"me.weldnor/swede/core/lang/swede/parser"
)

type Runner struct {
	registeredFuncs map[string]RegisteredFunc
	featureName     string
	scenarios       []Scenario
}

type RegisteredFunc func(*ScenarioContext) ScenarioExecutionResult

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

func (e *Runner) Run() {
	fmt.Printf("start running %s\n", e.featureName)

	var wg sync.WaitGroup

	ioMux := sync.Mutex{}

	for _, scenario := range e.scenarios {
		wg.Add(1)

		go func(scenario Scenario) {
			defer wg.Done()

			result := e.executeScenario(scenario)

			ioMux.Lock()
			defer ioMux.Unlock()

			if result.Status {
				fmt.Printf("[PASS] %s\n", scenario.Name)
			} else {
				fmt.Printf("[FAIL] %s: %s\n", scenario.Name, result.Message)
			}
		}(scenario)
	}

	wg.Wait()
}

func (e *Runner) executeScenario(scenario Scenario) ScenarioExecutionResult {
	scenarioContext := newScenarioContext()

	for _, step := range scenario.Steps {
		stepName := step.Name
		registeredFunc, ok := e.registeredFuncs[stepName]

		if !ok {
			return ScenarioExecutionResult{false, "step with name " + stepName + " is not registered"}
		}

		stepResult := registeredFunc(scenarioContext)

		if !stepResult.Status {
			return ScenarioExecutionResult{false, "failed on step " + stepName}
		}
	}

	return ScenarioExecutionResult{Status: true}
}
