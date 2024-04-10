package parser

import "me.weldnor/swede/core/lang/common"

const (
	ROOT       common.NodeType = "ROOT"
	UNEXPECTED common.NodeType = "UNEXPECTED"
	COMMENT    common.NodeType = "COMMENT"
	SCENARIO   common.NodeType = "SCENARIO"
	FEATURE    common.NodeType = "FEATURE"
	STEP       common.NodeType = "STEP"
	TAG        common.NodeType = "TAG"
)
