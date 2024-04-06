package parser

type Parser struct {
	source string
}

func NewParser(source string) *Parser {
	return &Parser{source: source}
}

// swede:step Add <first:type> and <second:type>

func (p *Parser) parse() {

}
