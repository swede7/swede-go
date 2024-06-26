package parser

import (
	"os"
	"strings"

	"github.com/swede7/swede-go/core/lang/common"
	"github.com/swede7/swede-go/core/lang/swede/lexer"
)

type Parser struct {
	lexemes []common.Lexeme
	nodes   []*common.Node
	errors  []common.ParserError
	pos     int
}

func ParseFile(path string) common.ParserResult {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic("can't read file " + path)
	}

	code := string(bytes)

	return ParseCode(code)
}

func ParseCode(code string) common.ParserResult {
	lexer := lexer.NewLexer(code)
	lexemes := lexer.Scan()

	parser := new(Parser)
	parser.lexemes = lexemes

	return parser.parse()
}

func (p *Parser) parse() common.ParserResult {
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

	rootNode := common.Node{}
	rootNode.Type = ROOT

	if len(p.nodes) > 0 {
		rootNode.StartPosition = p.nodes[0].StartPosition
		rootNode.EndPosition = p.nodes[len(p.nodes)-1].EndPosition
	}

	for _, node := range p.nodes {
		rootNode.AppendChild(node)
		rootNode.EndPosition = node.EndPosition
	}

	parserResult := common.ParserResult{}
	parserResult.Errors = p.errors
	parserResult.RootNode = rootNode

	return parserResult
}

func (p *Parser) peekLexeme() common.Lexeme {
	return p.lexemes[p.pos]
}

func (p *Parser) advance(count int) {
	p.pos += count
}

func (p *Parser) isEof() bool {
	return p.pos >= len(p.lexemes)
}

func (p *Parser) lookup(count int) common.Lexeme {
	return p.lexemes[p.pos+count]
}

func (p *Parser) lexemesLeft() int {
	return len(p.lexemes) - p.pos
}

func (p *Parser) addNode(node common.Node) {
	p.nodes = append(p.nodes, &node)
}

func (p *Parser) addError(startPosition common.Position, endPosition common.Position, message string) {
	p.errors = append(p.errors, common.ParserError{StartPosition: startPosition, EndPosition: endPosition, Message: message})
}

func (p *Parser) getPreviousLexeme() common.Lexeme {
	return p.lexemes[p.pos-1]
}

func (p *Parser) removeNodeByIndex(index int) {
	p.nodes = append(p.nodes[:index], p.nodes[index+1:]...)
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

func tagRule(p *Parser) bool {
	if p.isEof() {
		return false
	}

	if p.peekLexeme().Type != lexer.AT_CHR {
		return false
	}

	if p.lexemesLeft() <= 1 {
		return false
	}

	if p.lookup(1).Type != lexer.WORD {
		return false
	}

	atLexeme := p.peekLexeme()
	wordLexeme := p.lookup(1)

	tagNode := common.Node{
		Type:          TAG,
		StartPosition: atLexeme.StartPosition,
		EndPosition:   wordLexeme.EndPosition,
		Value:         wordLexeme.Value,
	}

	p.addNode(tagNode)
	p.advance(2)

	return true
}

func commentRule(p *Parser) bool {
	if p.isEof() {
		return false
	}

	if p.peekLexeme().Type != lexer.HASH_CHR {
		return false
	}

	hashLexeme := p.peekLexeme()

	sb := strings.Builder{}
	p.advance(1)

	for !p.isEof() && p.peekLexeme().Type != lexer.NL {
		currentLexeme := p.peekLexeme()
		sb.WriteString(currentLexeme.Value)
		p.advance(1)
	}

	commentNode := common.Node{
		Type:          COMMENT,
		StartPosition: hashLexeme.StartPosition,
		EndPosition:   p.peekLexeme().EndPosition,
		Value:         sb.String(),
	}

	p.addNode(commentNode)
	return true
}

func featureRule(p *Parser) bool {
	if p.isEof() {
		return false
	}

	if p.peekLexeme().Type != lexer.FEATURE_WORD {
		return false
	}
	featureWordLexeme := p.peekLexeme()

	var sb strings.Builder
	p.advance(1)

	for !p.isEof() && p.peekLexeme().Type != lexer.NL {
		currentLexeme := p.peekLexeme()
		sb.WriteString(currentLexeme.Value)
		p.advance(1)
	}

	featureNode := common.Node{Type: FEATURE, StartPosition: featureWordLexeme.StartPosition, EndPosition: p.getPreviousLexeme().EndPosition, Value: sb.String()}
	p.addNode(featureNode)
	return true
}

func addTagToFeatureRule(p *Parser) bool {
	if len(p.nodes) < 2 {
		return false
	}

	previousNode := p.nodes[len(p.nodes)-2]
	currentNode := p.nodes[len(p.nodes)-1]

	if previousNode.Type == TAG && currentNode.Type == FEATURE {
		currentNode.PrependChild(previousNode)
		p.removeNodeByIndex(len(p.nodes) - 2)
		return true
	}

	return false
}

func scenarioRule(p *Parser) bool {
	if p.isEof() {
		return false
	}

	if p.peekLexeme().Type != lexer.SCENARIO_WORD {
		return false
	}

	scenarioWordLexeme := p.peekLexeme()

	sb := strings.Builder{}

	p.advance(1)

	for !p.isEof() && p.peekLexeme().Type != lexer.NL {
		currentLexeme := p.peekLexeme()
		sb.WriteString(currentLexeme.Value)
		p.advance(1)
	}

	featureNode := common.Node{
		Type:          SCENARIO,
		StartPosition: scenarioWordLexeme.StartPosition,
		EndPosition:   p.getPreviousLexeme().EndPosition,
		Value:         sb.String(),
	}

	p.addNode(featureNode)
	return true
}

func addTagToScenarioRule(p *Parser) bool {
	if len(p.nodes) < 2 {
		return false
	}

	previousNode := p.nodes[len(p.nodes)-2]
	currentNode := p.nodes[len(p.nodes)-1]

	if previousNode.Type == TAG && currentNode.Type == SCENARIO {
		currentNode.PrependChild(previousNode)
		p.removeNodeByIndex(len(p.nodes) - 2)
		return true
	}
	return false
}

func stepRule(p *Parser) bool {
	if p.isEof() {
		return false
	}

	if p.peekLexeme().Type != lexer.DASH_CHR {
		return false
	}

	dashLexeme := p.peekLexeme()

	sb := strings.Builder{}
	p.advance(1)

	for !p.isEof() && p.peekLexeme().Type != lexer.NL {
		currentLexeme := p.peekLexeme()
		sb.WriteString(currentLexeme.Value)
		p.advance(1)
	}

	stepNode := common.Node{Type: STEP, StartPosition: dashLexeme.StartPosition, EndPosition: p.getPreviousLexeme().EndPosition, Value: sb.String()}
	p.addNode(stepNode)
	return true
}

func addStepToScenarioRule(p *Parser) bool {
	if len(p.nodes) < 2 {
		return false
	}

	previousNode := p.nodes[len(p.nodes)-2]
	currentNode := p.nodes[len(p.nodes)-1]

	if previousNode.Type == SCENARIO && currentNode.Type == STEP {
		previousNode.AppendChild(currentNode)
		p.removeNodeByIndex(len(p.nodes) - 1)
		return true
	}
	return false
}

func handleUnexpectedLexemeRule(p *Parser) bool {
	if p.isEof() {
		return false
	}

	unexpectedLexeme := p.peekLexeme()
	unexpectedNode := common.Node{Type: UNEXPECTED, StartPosition: unexpectedLexeme.StartPosition, EndPosition: unexpectedLexeme.EndPosition, Value: unexpectedLexeme.Value}

	p.addNode(unexpectedNode)
	p.addError(unexpectedLexeme.StartPosition, unexpectedLexeme.EndPosition, "unexpected lexeme")
	p.advance(1)
	return true
}

var validNodeTypes = map[common.NodeType]bool{
	UNEXPECTED: true,
	COMMENT:    true,
	FEATURE:    true,
	SCENARIO:   true,
}

func handleUnexpectedNodesRule(p *Parser) bool {
	if !p.isEof() {
		return false
	}

	someNodesWasProcessed := false
	for i, node := range p.nodes {
		if _, ok := validNodeTypes[node.Type]; !ok {
			wrapperNode := common.Node{Type: UNEXPECTED, StartPosition: node.StartPosition, EndPosition: node.EndPosition, Value: node.Value}

			wrapperNode.AppendChild(node)
			p.nodes[i] = &wrapperNode
			p.addError(node.StartPosition, node.EndPosition, "unexpected node")

			someNodesWasProcessed = true
		}
	}

	return someNodesWasProcessed
}
