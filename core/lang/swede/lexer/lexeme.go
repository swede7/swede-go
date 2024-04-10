package lexer

import "me.weldnor/swede/core/lang/common"

const (
	NL            common.LexemeType = "NL"
	AT_CHR        common.LexemeType = "AT_CHR"
	DASH_CHR      common.LexemeType = "DASH_CHR"
	HASH_CHR      common.LexemeType = "HASH_CHR"
	WORD          common.LexemeType = "WORD"
	SPACE         common.LexemeType = "SPACE"
	FEATURE_WORD  common.LexemeType = "FEATURE_WORD"
	SCENARIO_WORD common.LexemeType = "SCENARIO_WORD"
)
