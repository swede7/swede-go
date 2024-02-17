package formatter

import (
	"strings"

	"me.weldnor/swede/core/parser"
)

type Formatter struct {
	rootNode parser.Node
	sb       strings.Builder
}

func NewFormatter(rootNode *parser.Node) *Formatter {
	return &Formatter{
		rootNode: *rootNode,
	}
}

func (f Formatter) Format() (string, error) {
	for _, node := range f.rootNode.Children {
		f.formatNode(node)
	}

	return f.sb.String(), nil
}

func (f *Formatter) formatNode(node *parser.Node) {
	switch node.Type {
	case parser.FEATURE:
		f.formatFeature(node)
	case parser.COMMENT:
		f.formatComment(node)
	case parser.SCENARIO:
		f.formatScenario(node)
	}
}

func (f *Formatter) formatFeature(node *parser.Node) {
	tagNodes := node.GetChildrenByType(parser.TAG)

	for _, tagNode := range tagNodes {
		f.formatTag(tagNode)
	}
	f.sb.WriteString("\n")

	f.sb.WriteString("Feature: ")
	f.sb.WriteString(strings.TrimSpace(node.Value))
	f.sb.WriteString("\n\n")

}

func (f *Formatter) formatTag(node *parser.Node) {
	f.sb.WriteString("@")
	f.sb.WriteString(strings.TrimSpace(node.Value))
	f.sb.WriteString(" ")
}

func (f *Formatter) formatScenario(node *parser.Node) {
	tagNodes := node.GetChildrenByType(parser.TAG)

	for _, tagNode := range tagNodes {
		f.formatTag(tagNode)
	}
	f.sb.WriteString("\n")

	f.sb.WriteString("Scenario: ")
	f.sb.WriteString(strings.TrimSpace(node.Value))
	f.sb.WriteString("\n")

	stepNodes := node.GetChildrenByType(parser.STEP)

	for _, stepNode := range stepNodes {
		f.formatStep(stepNode)
	}
	f.sb.WriteString("\n")

}

func (f *Formatter) formatStep(node *parser.Node) {
	f.sb.WriteString("- ")
	f.sb.WriteString(strings.TrimSpace(node.Value))
	f.sb.WriteString("\n")
}

func (f *Formatter) formatComment(node *parser.Node) {
	f.sb.WriteString("# ")
	f.sb.WriteString(strings.TrimSpace(node.Value))
	f.sb.WriteString("\n\n")
}
