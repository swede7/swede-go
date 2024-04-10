package lexer

import (
	"errors"
	"strings"

	"me.weldnor/swede/core/lang/common"
)

func Lex(source string) ([]common.Lexeme, error) {
	_lexer, err := newLexer(source)
	if err != nil {
		return nil, err
	}

	return _lexer.lex(), nil
}

type lexer struct {
	source  string
	pos     int
	lexemes []common.Lexeme
}

const (
	TEXT      common.LexemeType = "text"
	L_BRACKET common.LexemeType = "l_bracket"
	R_BRACKET common.LexemeType = "r_bracket"
	COLON     common.LexemeType = "colon"
)

func newLexer(source string) (*lexer, error) {
	if strings.Contains(source, "\n") {
		return nil, errors.New("Lexer does not support multiline strings")
	}

	lexer := &lexer{}
	lexer.source = source

	return lexer, nil
}

func (l *lexer) lex() []common.Lexeme {
	for !l.isEnd() {
		c := l.curr()

		switch c {
		case '<':
			l.addLeftBracketLexeme()
		case '>':
			l.addRightBracketLexeme()
		case ':':
			l.addColonLexeme()
		default:
			l.lexText()
		}
	}

	return l.lexemes
}

func (l *lexer) lexText() {
	sb := strings.Builder{}
	startPos := l.getPosition()

	for !l.isEnd() {
		c := l.curr()

		if c == '>' || c == '<' || c == ':' {
			break
		}

		sb.WriteByte(c)
		l.next()
	}

	lexeme := common.Lexeme{
		Type:          TEXT,
		StartPosition: startPos,
		EndPosition:   l.getPosition().Inc(),
		Value:         sb.String(),
	}

	l.addLexeme(lexeme)
}

func (l *lexer) addLeftBracketLexeme() {
	lexeme := common.Lexeme{Type: L_BRACKET, StartPosition: l.getPosition(), EndPosition: l.getPosition(), Value: "<"}
	l.addLexeme(lexeme)
	l.next()
}

func (l *lexer) addRightBracketLexeme() {
	lexeme := common.Lexeme{Type: R_BRACKET, StartPosition: l.getPosition(), EndPosition: l.getPosition(), Value: ">"}
	l.addLexeme(lexeme)
	l.next()
}

func (l *lexer) addColonLexeme() {
	lexeme := common.Lexeme{Type: COLON, StartPosition: l.getPosition(), EndPosition: l.getPosition(), Value: ":"}
	l.addLexeme(lexeme)
	l.next()
}

func (l *lexer) getPreviousLexeme() *common.Lexeme {
	if len(l.lexemes) == 0 {
		return nil
	}

	return &l.lexemes[len(l.lexemes)-1]
}

func (l *lexer) addLexeme(lexeme common.Lexeme) {
	l.lexemes = append(l.lexemes, lexeme)
}

func (l *lexer) curr() byte {
	return l.source[l.pos]
}

func (l *lexer) next() {
	l.pos++
}

func (l *lexer) getPosition() common.Position {
	return common.Position{Line: 0, Column: l.pos, Offset: l.pos}
}

func (l *lexer) isEnd() bool {
	return l.pos >= len(l.source)
}
