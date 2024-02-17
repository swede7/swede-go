package parser

import (
	"me.weldnor/swede/core/common"
)

type Node struct {
	StartPosition common.Position
	EndPosition   common.Position
	Value         string
	Type          NodeType
	Children      []*Node
}

type NodeType string

const (
	ROOT       NodeType = "ROOT"
	UNEXPECTED NodeType = "UNEXPECTED"
	COMMENT    NodeType = "COMMENT"
	SCENARIO   NodeType = "SCENARIO"
	FEATURE    NodeType = "FEATURE"
	STEP       NodeType = "STEP"
	TAG        NodeType = "TAG"
)

func (node *Node) AppendChild(child *Node) {
	node.Children = append(node.Children, child)
}

func (node *Node) PrependChild(child *Node) {
	node.Children = append([]*Node{child}, node.Children...)
}
