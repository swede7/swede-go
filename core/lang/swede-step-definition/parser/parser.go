package parser

import (
	"me.weldnor/swede/core/lang/common"
	"me.weldnor/swede/core/lang/swede-step-definition/lexer"
	"me.weldnor/swede/core/lang/swede-step-definition/model"
	"regexp"
)

type parser struct {
	source  string
	lexemes []common.Lexeme
	nodes   []common.Node
	pos     int

	rootNode common.Node //after Parse()
	model    model.StepDefinition
}

type ParserResult struct {
	common.ParserResult
	StepDefinition model.StepDefinition
}

func Parse(source string) (*ParserResult, error) {
	lexemes, err := lexer.Lex(source)
	if err != nil {
		return nil, err
	}

	_parser := newParser(lexemes)
	_parser.source = source

	return _parser.Parse()
}

func (p *parser) Parse() (*ParserResult, error) {
	p.applyParseRules()

	_model := p.getModel()

	result := ParserResult{}
	result.StepDefinition = _model
	result.RootNode = p.rootNode
	result.Errors = make([]common.ParserError, 0) //todo fixme

	return &result, nil
}

func (p *parser) getModel() model.StepDefinition {
	_model := model.StepDefinition{}
	_model.Text = p.source

	regexString := ""

	variables := make([]model.Variable, 0)

	for _, node := range p.rootNode.Children {
		if node.Type == TEXT {
			regexString += regexp.QuoteMeta(node.Value)
		}

		if node.Type == VARIABLE {
			variableName := node.GetChildByType(VARIABLE_NAME).Value

			variableTypeAsString := node.GetChildByType(VARIABLE_TYPE).Value
			variableType, err := model.GetVariableTypeByName(variableTypeAsString)

			if err != nil {
				panic(err) //todo fixme
			}

			variables = append(variables, model.Variable{
				Name: variableName,
				Type: variableType,
			})

			regexString += variableType.RegexTemplate()
		}
	}

	_model.Variables = variables
	_model.Regex = regexp.MustCompile(regexString)
	return _model
}

func (p *parser) applyParseRules() {
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

	p.rootNode = p.nodes[0]
}

// region Node types

const (
	ROOT          common.NodeType = "root"
	TEXT          common.NodeType = "text"
	VARIABLE      common.NodeType = "variable"
	VARIABLE_NAME common.NodeType = "variable_name"
	VARIABLE_TYPE common.NodeType = "variable_type"

	UNEXPECTED common.NodeType = "unexpected"

	L_BRACKET common.NodeType = "l_bracket"
	R_BRACKET common.NodeType = "r_bracket"
	COLON     common.NodeType = "colon"
)

// region Parse rules

type parseRule func(*parser) bool

var parseRules = []parseRule{
	variableRule,
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

	newNode.StartPosition = currentLexeme.StartPosition
	newNode.EndPosition = currentLexeme.EndPosition
	newNode.Value = currentLexeme.Value

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

	for i := 0; i < len(p.nodes)-1; i++ {
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

	p.deleteNodesFromEnd(5)
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
		rootNode.StartPosition = common.Position{}
		rootNode.EndPosition = common.Position{}

		p.nodes = append(p.nodes, rootNode)
		return true
	}

	rootNode.StartPosition = p.nodes[0].StartPosition
	rootNode.EndPosition = p.nodes[0].EndPosition

	for i := 0; i < len(p.nodes); i++ {
		rootNode.AppendChild(&p.nodes[i])
		rootNode.EndPosition = p.nodes[i].EndPosition
	}

	p.nodes = []common.Node{rootNode}
	return true
}

func (p *parser) getCurrentLexeme() common.Lexeme {
	return p.lexemes[p.pos]
}

func newParser(lexemes []common.Lexeme) *parser {
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

func (p *parser) deleteNodesFromEnd(count int) {
	p.nodes = p.nodes[:len(p.nodes)-count]
}

func (p *parser) isEnd() bool {
	return p.pos >= len(p.lexemes)
}
