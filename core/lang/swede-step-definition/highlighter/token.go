package highlight

import (
	"github.com/swede7/swede-go/core/lang/common"
)

type TokenType string

const (
	Variable TokenType = "variable"
	Type     TokenType = "type"
)

type Token struct {
	Start     common.Position
	Length    int
	Type      TokenType
	modifiers []string
}
