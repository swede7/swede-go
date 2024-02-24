package runner

type Feature struct {
	Name      string
	Scenarios []Scenario
}

type Scenario struct {
	Name  string
	Steps []Step
}

type Step struct {
	Name string
}

type ScenarioExecutionResult struct {
	Status  bool
	Message string
}
