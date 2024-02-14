package core

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
}

func scanSpace(l *Lexer) bool {

}

func scanHash(l *Lexer) bool {

}

func scanDash(l *Lexer) bool {

}

func scanAt(l *Lexer) bool {

}

func scanNl(l *Lexer) bool {

}
