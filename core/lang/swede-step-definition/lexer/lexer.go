package lexer

import "strings"

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
	Type  LexemeType
	Start int
	End   int
	Text  string
}

type LexemeType string

const (
	TEXT      LexemeType = "text"
	L_BRACKET LexemeType = "l_bracket"
	R_BRACKET LexemeType = "r_bracket"
	COLON     LexemeType = "colon"
)

func newLexer(source string) *lexer {
	lexer := &lexer{}
	lexer.source = source

	return lexer
}

func (l *lexer) lex() []Lexeme {
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
	startPos := l.pos

	for !l.isEnd() {
		c := l.curr()

		if c == '>' || c == '<' || c == ':' {
			break
		}

		sb.WriteByte(c)
		l.next()
	}

	lexeme := Lexeme{
		Type:  TEXT,
		Start: startPos,
		End:   l.pos - 1,
		Text:  sb.String(),
	}

	l.addLexeme(lexeme)
}

func (l *lexer) addLeftBracketLexeme() {
	lexeme := Lexeme{Type: L_BRACKET, Start: l.pos, End: l.pos, Text: "<"}
	l.addLexeme(lexeme)
	l.next()
}

func (l *lexer) addRightBracketLexeme() {
	lexeme := Lexeme{Type: R_BRACKET, Start: l.pos, End: l.pos, Text: ">"}
	l.addLexeme(lexeme)
	l.next()
}

func (l *lexer) addColonLexeme() {
	lexeme := Lexeme{Type: COLON, Start: l.pos, End: l.pos, Text: ":"}
	l.addLexeme(lexeme)
	l.next()
}

func (l *lexer) getPreviousLexeme() *Lexeme {
	if len(l.lexemes) == 0 {
		return nil
	}

	return &l.lexemes[len(l.lexemes)-1]
}

func (l *lexer) addLexeme(lexeme Lexeme) {
	l.lexemes = append(l.lexemes, lexeme)
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
