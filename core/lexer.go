package core

import (
	"strings"
	"unicode"
)

type Lexer struct {
	source  string
	offset  int
	line    int
	column  int
	lexemes []Lexeme
}

func (l *Lexer) isAtEnd() bool {
	return l.offset >= len(l.source)
}

func (l *Lexer) advance() {
	l.offset++
	l.column++
}

func (l *Lexer) charsLeft() int {
	return len(l.source) - l.offset
}

func (l *Lexer) peek() uint8 {
	if l.isAtEnd() {
		return '\000'
	}
	return l.source[l.offset]
}

func (l *Lexer) matchChar(expected uint8) bool {
	if l.isAtEnd() {
		return false
	}
	if l.peek() != expected {
		return false
	}
	l.advance()
	return true
}

func (l *Lexer) matchString(expected string) bool {
	if l.isAtEnd() {
		return true
	}

	if len(expected) > l.charsLeft() {
		return false
	}

	startPosition := l.getPosition()

	for i := 0; i < len(expected); i++ {
		if expected[i] != l.peek() {
			//rollback
			l.setPosition(startPosition)
			return false
		}
		l.advance()
	}
	return true
}

func (l *Lexer) addToken(lexemeType LexemeType, startPosition Position, endPosition Position, value string) {
	l.lexemes = append(l.lexemes, Lexeme{lexemeType, startPosition, endPosition, value})
}

func (l *Lexer) getPosition() Position {
	return Position{offset: l.offset, line: l.line, column: l.column}
}

func (l *Lexer) setPosition(position Position) {
	l.offset = position.offset
	l.line = position.line
	l.column = position.column
}

const (
	featureWord  string = "Feature:"
	scenarioWord string = "Feature:"
)

func NewLexer(source string) *Lexer {
	return &Lexer{
		source: source,
	}
}

func (l *Lexer) Scan() []Lexeme {
	for !l.isAtEnd() {
		l.scanNextToken()
	}
	return l.lexemes
}

func (l *Lexer) scanNextToken() {
	for _, scanFunction := range scanFunctions {
		if scanFunction(l) {
			return
		}
	}

	if !l.isAtEnd() {
		panic("oops! can't read all tokens")
	}
}

type scanFunction = func(lexer *Lexer) bool

var scanFunctions = []scanFunction{
	scanNl,
	scanAt,
	scanDash,
	scanHash,
	scanSpace,
	scanWord,
}

func scanWord(l *Lexer) bool {
	if l.isAtEnd() {
		return false
	}

	startPosition := l.getPosition()

	if l.matchString(scenarioWord) {
		l.addToken(SCENARIO_WORD, startPosition, Position{l.offset - 1, l.line, l.column - 1}, scenarioWord)
	}

	if l.matchString(featureWord) {
		l.addToken(FEATURE_WORD, startPosition, Position{l.offset - 1, l.line, l.column - 1}, featureWord)
	}

	sb := strings.Builder{}

	for !l.isAtEnd() {
		currentChar := l.peek()
		if unicode.IsSpace(rune(currentChar)) {
			break
		}
		sb.WriteRune(rune(currentChar))
		l.advance()
	}

	l.addToken(WORD, startPosition, Position{l.offset - 1, l.line, l.column - 1}, sb.String())
	return true
}

func scanSpace(l *Lexer) bool {
	if l.isAtEnd() {
		return false
	}

	startPosition := l.getPosition()

	sb := strings.Builder{}

	for !l.isAtEnd() && l.peek() == ' ' || l.peek() == '\t' {
		sb.WriteByte(l.peek())
		l.advance()
	}

	if sb.Len() == 0 {
		return false
	}

	l.addToken(SPACE, startPosition, Position{l.offset - 1, l.line, l.column - 1}, sb.String())
	return true
}

func scanHash(l *Lexer) bool {
	if l.isAtEnd() {
		return false
	}

	startPosition := l.getPosition()

	if !l.matchChar('#') {
		return false
	}

	l.addToken(HASH_CHR, startPosition, startPosition, "#")
	return true
}

func scanDash(l *Lexer) bool {
	if l.isAtEnd() {
		return false
	}

	startPosition := l.getPosition()

	if !l.matchChar('-') {
		return false
	}

	l.addToken(DASH_CHR, startPosition, startPosition, "-")
	return true
}

func scanAt(l *Lexer) bool {
	if l.isAtEnd() {
		return false
	}

	startPosition := l.getPosition()

	if !l.matchChar('@') {
		return false
	}

	l.addToken(AT_CHR, startPosition, startPosition, "@")
	return true
}

func scanNl(l *Lexer) bool {
	if l.isAtEnd() {
		return false
	}

	startPosition := l.getPosition()

	if l.matchString("\r\n") {
		l.addToken(NL, startPosition, Position{l.offset - 1, l.line, l.column - 1}, "\r\n")
		l.line++
		l.column = 0
		return true
	}

	if l.matchString("\n") {
		l.addToken(NL, startPosition, Position{l.offset - 1, l.line, l.column - 1}, "\n")
		l.line++
		l.column = 0
		return true
	}
	return false
}
