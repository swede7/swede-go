package lexer

import (
	"strings"
	"unicode"
)

func Lex(source string) []Lexeme {
	_lexer := newLexer(source)

	return _lexer.lex()
}

type lexer struct {
	source  string
	pos     int
	lexemes []Lexeme
}

type Lexeme struct {
	Start int
	End   int
	Type  LexemeType
	Text  string
}

type LexemeType string

const (
	TEXT     LexemeType = "text"
	VARIABLE LexemeType = "variable"
)

func newLexer(source string) *lexer {
	lexer := &lexer{}
	lexer.source = source

	return lexer
}

// swede:step Add <first:int> and <second:int>

func (l *lexer) lex() []Lexeme {
	for !l.isEnd() {
		c := l.curr()

		if c == '<' {
			l.next()
			l.lexVariable()
			continue
		}

		if isTextCharacter(c) {
			l.lexText()
			continue
		}
	}

	return l.lexemes
}

func (l *lexer) lexText() {
	if l.isEnd() {
		return
	}

	startPos := l.pos

	sb := strings.Builder{}

	//todo possible bug

	for !l.isEnd() {
		c := l.curr()

		if c == '<' {
			break
		}

		sb.WriteByte(l.curr())
		l.next()
	}

	endPos := l.pos - 1

	lexeme := Lexeme{
		Type:  TEXT,
		Start: startPos,
		End:   endPos,
		Text:  sb.String(),
	}

	l.lexemes = append(l.lexemes, lexeme)
}

func (l *lexer) lexVariable() {
	if l.isEnd() {
		return
	}

	startPos := l.pos

	sb := strings.Builder{}

	//todo possible bug

	for !l.isEnd() {
		c := l.curr()

		if c == '>' {
			l.next() //skip right brace

			break
		}

		sb.WriteByte(l.curr())
		l.next()
	}

	endPos := l.pos - 1

	lexeme := Lexeme{
		Type:  TEXT,
		Start: startPos,
		End:   endPos,
	}

	l.lexemes = append(l.lexemes, lexeme)
}

func isTextCharacter(c byte) bool {
	return unicode.IsLetter(rune(c))
}

func (l *lexer) curr() byte {
	return l.source[l.pos]
}

func (l *lexer) next() {
	l.pos++
}

func (l *lexer) isEnd() bool {
	return l.pos >= len(l.source)
}
