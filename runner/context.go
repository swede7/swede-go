package runner

type ScenarioContext struct {
	Variables map[string]any
}

func newScenarioContext() *ScenarioContext {
	return &ScenarioContext{Variables: make(map[string]any)}
}
