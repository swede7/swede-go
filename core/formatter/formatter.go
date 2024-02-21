package formatter

import (
	"strings"
	"sync"

	"me.weldnor/swede/core/parser"
)

type Formatter struct {
	rootNode parser.Node
}

func NewFormatter(rootNode *parser.Node) *Formatter {
	return &Formatter{
		rootNode: *rootNode,
	}
}

func (f *Formatter) FormatParallel() (string, error) {
	var wg sync.WaitGroup
	results := make([]string, len(f.rootNode.Children))

	for i, node := range f.rootNode.Children {
		wg.Add(1)
		go func(i int, node *parser.Node) {
			defer wg.Done()
			results[i] = f.formatNode(node)
		}(i, node)
	}
	wg.Wait()

	return strings.Join(results, ""), nil
}

func (f *Formatter) formatNode(node *parser.Node) string {
	var sb strings.Builder

	switch node.Type {
	case parser.FEATURE:
		f.formatFeature(&sb, node)
	case parser.COMMENT:
		f.formatComment(&sb, node)
	case parser.SCENARIO:
		f.formatScenario(&sb, node)
	}

	return sb.String()
}

// ... rest of the Formatter methods, modified to take a *strings.Builder as the first argument ...
func (f *Formatter) formatFeature(sb *strings.Builder, node *parser.Node) {
	tagNodes := node.GetChildrenByType(parser.TAG)

	for _, tagNode := range tagNodes {
		f.formatTag(sb, tagNode)
	}
	sb.WriteString("\n")

	sb.WriteString("Feature: ")
	sb.WriteString(strings.TrimSpace(node.Value))
	sb.WriteString("\n\n")

}

func (f *Formatter) formatTag(sb *strings.Builder, node *parser.Node) {
	sb.WriteString("@")
	sb.WriteString(strings.TrimSpace(node.Value))
	sb.WriteString(" ")
}

func (f *Formatter) formatScenario(sb *strings.Builder, node *parser.Node) {
	tagNodes := node.GetChildrenByType(parser.TAG)

	for _, tagNode := range tagNodes {
		f.formatTag(sb, tagNode)
	}
	sb.WriteString("\n")

	sb.WriteString("Scenario: ")
	sb.WriteString(strings.TrimSpace(node.Value))
	sb.WriteString("\n")

	stepNodes := node.GetChildrenByType(parser.STEP)

	for _, stepNode := range stepNodes {
		f.formatStep(sb, stepNode)
	}
	sb.WriteString("\n")

}

func (f *Formatter) formatStep(sb *strings.Builder, node *parser.Node) {
	sb.WriteString("- ")
	sb.WriteString(strings.TrimSpace(node.Value))
	sb.WriteString("\n")
}

func (f *Formatter) formatComment(sb *strings.Builder, node *parser.Node) {
	sb.WriteString("# ")
	sb.WriteString(strings.TrimSpace(node.Value))
	sb.WriteString("\n\n")
}
