package common

type Position struct {
	Offset int
	Line   int
	Column int
}

func (p Position) Inc() Position {
	if p.Column == 0 {
		panic("not supported yet")
	}

	return Position{
		Line:   p.Line,
		Offset: p.Offset - 1,
		Column: p.Column - 1,
	}
}
