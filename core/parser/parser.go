package parser

import (
	"me.weldnor/swede/core/lexer"
)

type Parser struct {
	position int
	lexemes  []lexer.Lexeme
	nodes    []Node
	errors   []Error
	pos      int
}

func (p *Parser) NewParser(lexemes []lexer.Lexeme) *Parser {
	return &Parser{
		position: 0,
		lexemes:  lexemes,
	}
}

func (p *Parser) peekLexeme() lexer.Lexeme {
	return p.lexemes[p.pos]
}

func (p *Parser) advance(count int) {
	p.pos += count
}

func (p *Parser) isEof() bool {
	return p.pos >= len(p.lexemes)
}

func (p *Parser) lookup(count int) lexer.Lexeme {
	return p.lexemes[p.pos+count]
}

func (p *Parser) lexemesLeft() int {
	return len(p.lexemes) - p.pos
}

func (p *Parser) addNode(node Node) {
	p.nodes = append(p.nodes, node)
}

func (p *Parser) getPreviousLexeme() lexer.Lexeme {
	return p.lexemes[p.pos-1]
}

type parseRule func(*Parser) bool

var parseRules []parseRule = []parseRule{
	addTagToFeatureRule,
	addTagToScenarioRule,
	addStepToScenarioRule,
	skipSpacesAndNlRule,
	commentRule,
	tagRule,
	featureRule,
	scenarioRule,
	stepRule,
	handleUnexpectedLexemeRule,
	handleUnexpectedNodesRule,
}

func (p *Parser) Parse() {
	for {
		anyRuleWasApplied := false

		for rule := range rules {

		}

	}
}

func skipSpacesAndNlRule(p *Parser) bool {
	if p.isEof() {
		return false
	}

	if p.peekLexeme().Type == lexer.SPACE || p.peekLexeme().Type == lexer.NL {
		p.advance(1)
		return true
	}
	return false
}

func