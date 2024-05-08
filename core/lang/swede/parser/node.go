package parser

import "github.com/swede7/swede-go/core/lang/common"

const (
	ROOT       common.NodeType = "ROOT"
	UNEXPECTED common.NodeType = "UNEXPECTED"
	COMMENT    common.NodeType = "COMMENT"
	SCENARIO   common.NodeType = "SCENARIO"
	FEATURE    common.NodeType = "FEATURE"
	STEP       common.NodeType = "STEP"
	TAG        common.NodeType = "TAG"
)
