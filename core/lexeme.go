package core

type Lexeme struct {
	Type          LexemeType
	startPosition Position
	endPosition   Position
	value         string
}

type LexemeType string

const (
	NL            LexemeType = "NL"
	AT_CHR        LexemeType = "AT_CHR"
	DASH_CHR      LexemeType = "DASH_CHR"
	HASH_CHR      LexemeType = "HASH_CHR"
	WORD          LexemeType = "WORD"
	SPACE         LexemeType = "SPACE"
	FEATURE_WORD  LexemeType = "FEATURE_WORD"
	SCENARIO_WORD LexemeType = "SCENARIO_WORD"
)
