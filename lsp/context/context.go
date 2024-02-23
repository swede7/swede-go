package context

var instance *LspContext

type LspContext struct {
	Code string
}

func GetContext() *LspContext {
	if instance == nil {
		instance = &LspContext{}
	}

	return instance
}
