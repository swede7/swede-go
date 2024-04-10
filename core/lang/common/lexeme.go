package common

type Lexeme struct {
	Type          LexemeType
	StartPosition Position
	EndPosition   Position
	Value         string
}

type LexemeType string
