package runner

type Context struct {
	variables map[string]any
}

func Clone() *Context {
	return &Context{variables: make(map[string]any)}
}

func (c *Context) GetIntVariable(name string) int {
	return c.GetVariable(name).(int)
}

func (c *Context) GetStringVariable(name string) string {
	return c.GetVariable(name).(string)
}

func (c *Context) GetVariable(name string) any {
	return c.variables[name]
}

func (c *Context) SetVariable(name string, value any) {
	c.variables[name] = value
}
