package parser

import "me.weldnor/swede/core/common"

type Node struct {
	StartPosition common.Position
	EndPosition   common.Position
	Value         string
	Type          NodeType
}

type NodeType string

const (
	ROOT       NodeType = "ROOT"
	UNEXPECTED          = "UNEXPECTED"
	COMMENT             = "COMMENT"
	SCENARIO            = "SCENARIO"
	FEATURE             = "FEATURE"
	STEP                = "STEP"
	TAG                 = "TAG"
)
