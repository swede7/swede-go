package core

type Lexeme struct {
	Type          LexemeType
	startPosition Position
	endPosition   Position
	value         string
}

type LexemeType int

const (
	NL LexemeType = iota
	AT_CHR
	DASH_CHR
	HASH_CHR
	WORD
	SPACE
	FEATURE_WORD
	SCENARIO_WORD
)
