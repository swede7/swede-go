package runner

type Context struct {
	variables map[string]any
}

func newContext() *Context {
	return &Context{variables: make(map[string]any)}
}

func (c *Context) GetVariable(name string) any {
	return c.variables[name]
}

func (c *Context) SetVariable(name string, value any) {
	c.variables[name] = value
}
