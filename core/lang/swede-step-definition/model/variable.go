package model

import (
	"errors"
)

type Variable struct {
	Name string
	Type VariableType
}

type ParsedVariable struct {
	Variable
	Value any
}

type VariableType string

const (
	String VariableType = "string"
	Int    VariableType = "int"
)

var VariableTypes = map[VariableType]bool{
	String: true,
	Int:    true,
}

func (t VariableType) RegexTemplate() string {
	switch t {
	case String:
		return "\"([\\w ]*)\""
	case Int:
		return "(\\d+)"
	default:
		panic("incorrect variable type")
	}
}

func GetVariableTypeByName(name string) (VariableType, error) {
	if VariableTypes[VariableType(name)] {
		return VariableType(name), nil
	}

	return "", errors.New("variable type not found")
}
