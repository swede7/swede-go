package model

import (
	"regexp"
	"strconv"
)

type StepDefinition struct {
	Text      string
	Variables []Variable
	Regex     *regexp.Regexp
}

func (s *StepDefinition) Check(stepText string) bool {
	return s.Regex.MatchString(stepText)
}

func (s *StepDefinition) Equals(another StepDefinition) bool {
	return s.Regex == another.Regex
}

func (s *StepDefinition) GetParsedValues(stepText string) []ParsedVariable {
	parsedVariables := make([]ParsedVariable, 0)

	valueAsStrings := s.Regex.FindStringSubmatch(stepText)[1:]

	for i, variable := range s.Variables {
		valueAsString := valueAsStrings[i]
		convertedValue := convertStringByVariableType(valueAsString, variable.Type)
		parsedVariables = append(parsedVariables, ParsedVariable{
			Variable: variable,
			Value:    convertedValue,
		})
	}

	return parsedVariables
}

func convertStringByVariableType(valueAsString string, variableType VariableType) any {
	switch variableType {
	case String:
		return valueAsString
	case Int:
		{
			value, err := strconv.Atoi(valueAsString)
			if err != nil {
				panic(err)
			}
			return value
		}
	default:
		panic("Unsupported variable type")
	}
}
