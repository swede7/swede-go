package parser

import (
	"me.weldnor/swede/core/common"
	"me.weldnor/swede/core/lang/swede-step-definition/lexer"
)

type parser struct {
	source  string
	lexemes []lexer.Lexeme
	nodes   []common.Node
	pos     int
}

func Parse(source string) common.Node {
	lexemes := lexer.Lex(source)
	_parser := newParser(lexemes)
	return _parser.Parse()
}

func (p *parser) Parse() common.Node {
	for {
		anyRuleWasApplied := false

		for _, rule := range parseRules {
			if rule(p) {
				anyRuleWasApplied = true
				break
			}
		}

		if !anyRuleWasApplied {
			break
		}
	}

	return p.nodes[0]
}

// region Node types

const (
	ROOT          common.NodeType = "root"
	TEXT          common.NodeType = "variable"
	VARIABLE      common.NodeType = "variable"
	VARIABLE_NAME common.NodeType = "variable_name"
	VARIABLE_TYPE common.NodeType = "variable_type"

	UNEXPECTED common.NodeType = "unexpected"

	L_BRACKET common.NodeType = "l_bracket"
	R_BRACKET common.NodeType = "r_bracket"
	COLON     common.NodeType = "colon"
)

// region parse rules

type parseRule func(*parser) bool

var parseRules = []parseRule{
	VariableDefinitionRule,
	addNodeRule,
	mergeColonAndTextRule,
	mergeToRootNode,
}

func addNodeRule(p *parser) bool {
	if p.isEnd() {
		return false
	}

	newNode := common.Node{}
	currentLexeme := p.getCurrentLexeme()

	newNode.StartPosition = currentLexeme.Start
	newNode.EndPosition = currentLexeme.End
	newNode.Value = currentLexeme.Text

	switch currentLexeme.Type {
	case lexer.TEXT:
		newNode.Type = TEXT
	case lexer.COLON:
		newNode.Type = COLON
	case lexer.L_BRACKET:
		newNode.Type = L_BRACKET
	case lexer.R_BRACKET:
		newNode.Type = R_BRACKET

	default:
		panic("incorrect state")
	}

	p.addNode(newNode)

	p.next()

	return true
}

func mergeColonAndTextRule(p *parser) bool {
	changed := false

	updatedNodes := make([]common.Node, 0)

	for i := 0; i < len(p.nodes); i++ {
		currentNode := p.nodes[i]

		// last node
		if i+1 == p.pos {
			updatedNodes = append(updatedNodes, currentNode)
			continue
		}

		nextNode := p.nodes[i+1]
		// text, colon || colon, text -> text
		if currentNode.Type == TEXT && nextNode.Type == COLON || currentNode.Type == COLON && nextNode.Type == TEXT {
			newTextNode := common.Node{}
			newTextNode.Type = TEXT
			newTextNode.StartPosition = currentNode.StartPosition
			newTextNode.EndPosition = nextNode.EndPosition
			newTextNode.Value = currentNode.Value + newTextNode.Value

			updatedNodes = append(updatedNodes, newTextNode)
			changed = true
			i++ // skip next (already merged) node
		}
	}

	return changed
}

func variableRule(p *parser) bool {
	countOfNodes := len(p.nodes)

	if countOfNodes < 5 {
		return false
	}

	if p.getNodeFromEnd(0).Type != R_BRACKET {
		return false
	}

	if p.getNodeFromEnd(1).Type != TEXT {
		return false
	}

	if p.getNodeFromEnd(2).Type != COLON {
		return false
	}

	if p.getNodeFromEnd(3).Type != TEXT {
		return false
	}

	if p.getNodeFromEnd(4).Type != L_BRACKET {
		return false
	}

	variableNameNode := common.Node{}
	variableNameNode.Type = VARIABLE_NAME
	variableNameNode.StartPosition = p.getNodeFromEnd(3).StartPosition
	variableNameNode.EndPosition = p.getNodeFromEnd(3).EndPosition
	variableNameNode.Value = p.getNodeFromEnd(3).Value

	variableTypeNode := common.Node{}
	variableTypeNode.Type = VARIABLE_TYPE
	variableTypeNode.StartPosition = p.getNodeFromEnd(1).StartPosition
	variableTypeNode.EndPosition = p.getNodeFromEnd(1).EndPosition
	variableTypeNode.Value = p.getNodeFromEnd(1).Value

	variableNode := common.Node{}
	variableNode.StartPosition = p.getNodeFromEnd(4).StartPosition
	variableNode.EndPosition = p.getNodeFromEnd(0).EndPosition
	variableNode.Type = VARIABLE

	variableNode.AppendChild(&variableNameNode)
	variableNode.AppendChild(&variableTypeNode)

	p.addNode(variableNode)

	return true
}

// region Utils methods
func VariableDefinitionRule(parser *parser) bool {
	return true
}

func addColonToTextRule() {
}

func mergeToRootNode(p *parser) bool {
	if len(p.nodes) == 1 && p.nodes[0].Type == ROOT {
		return false
	}

	rootNode := common.Node{}
	rootNode.Type = ROOT

	if len(p.nodes) == 0 {
		rootNode.StartPosition = 0
		rootNode.EndPosition = 0

		p.nodes = append(p.nodes, rootNode)
		return true
	}

	if len(p.nodes) == 1 {
		rootNode.StartPosition = p.nodes[0].StartPosition
		rootNode.EndPosition = p.nodes[0].EndPosition
		return true
	}

	rootNode.StartPosition = 0
	rootNode.EndPosition = 0
	for i := 0; i < len(p.nodes); i++ {
		rootNode.AppendChild(&p.nodes[i])
	}

	p.nodes = []common.Node{rootNode}
	return true
}   

func (p *parser) getCurrentLexeme() lexer.Lexeme {
	return p.lexemes[p.pos]
}

func newParser(lexemes []lexer.Lexeme) *parser {
	return &parser{lexemes: lexemes}
}

func (p *parser) addNode(node common.Node) {
	p.nodes = append(p.nodes, node)
}

func (p *parser) getNodeFromEnd(offset int) *common.Node {
	countOfNodes := len(p.nodes)
	return &p.nodes[countOfNodes-offset-1]
}

func (p *parser) next() {
	p.pos++
}



func (p *parser) isEnd() bool {
	return p.pos >= len(p.lexemes)
}
