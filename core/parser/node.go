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

func (node *Node) GetChildrenByType(nodeType NodeType) []*Node {
	result := make([]*Node, 0)

	for _, child := range node.Children {
		if child.Type == nodeType {
			result = append(result, child)
		}
	}
	return result
}

func (node *Node) GetChildByType(nodeType NodeType) *Node {
	children := node.GetChildrenByType(nodeType)
	return children[0]
}

func VisitNode(node *Node, action func(*Node)) {
	action(node)

	for _, child := range node.Children {
		VisitNode(child, action)
	}
}
