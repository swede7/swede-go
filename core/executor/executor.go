package executor

import (
	"fmt"
	"sync"
)

type Executor struct {
	registeredFuncs map[string]RegisteredFunc
	scenarios       []Scenario
}

type StepExecutionResult struct {
	Status  bool
	Message string
}

type ScenarioExecutionResult struct {
	Name    string
	Status  bool
	Message string
}

type RegisteredFunc func(*ScenarioContext) StepExecutionResult

func (e *Executor) RegisterFunc(name string, registeredFunc RegisteredFunc) {
	if e.registeredFuncs == nil {
		e.registeredFuncs = make(map[string]RegisteredFunc)
	}
	e.registeredFuncs[name] = registeredFunc
}

func (e *Executor) ExecuteScenariosParallel() {
	var wg sync.WaitGroup
	results := make(chan ScenarioExecutionResult, len(e.scenarios))

	for _, scenario := range e.scenarios {
		wg.Add(1)
		go func(scenario Scenario) {
			defer wg.Done()
			result := e.ExecuteScenario(scenario)
			results <- result
		}(scenario)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Printf("Scenario %s: %s\n", result.Name, result.Message)
	}
}

func (e *Executor) ExecuteScenario(scenario Scenario) ScenarioExecutionResult {
	scenarioName := scenario.Name
	scenarioContext := &ScenarioContext{}

	for _, step := range scenario.Steps {
		stepName := step.Name
		registeredFunc, ok := e.registeredFuncs[stepName]

		if !ok {
			return ScenarioExecutionResult{scenarioName, false, "step with name " + stepName + " is not registered"}
		}

		stepResult := registeredFunc(scenarioContext)

		if !stepResult.Status {
			return ScenarioExecutionResult{scenarioName, false, "failed on step " + stepName}
		}
	}

	return ScenarioExecutionResult{scenarioName, true, "ok"}
}

type Scenario struct {
	Name  string
	Steps []Step
}

type Step struct {
	Name string
}
