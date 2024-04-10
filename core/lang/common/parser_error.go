package common

type ParserError struct {
	StartPosition Position
	EndPosition   Position
	Message       string
}
