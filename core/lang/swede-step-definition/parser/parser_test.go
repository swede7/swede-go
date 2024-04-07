package parser

import (
	"reflect"
	"testing"

	"me.weldnor/swede/core/common"
	"me.weldnor/swede/core/lang/swede-step-definition/lexer"
)

func TestParse(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name string
		args args
		want common.Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.source); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_Parse(t *testing.T) {
	tests := []struct {
		name string
		p    *parser
		want common.Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Parse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_variableDefinitionRule(t *testing.T) {
	type args struct {
		parser *parser
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := variableRule(tt.args.parser); got != tt.want {
				t.Errorf("variableDefinitionRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariableDefinitionRule(t *testing.T) {
	type args struct {
		parser *parser
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VariableDefinitionRule(tt.args.parser); got != tt.want {
				t.Errorf("VariableDefinitionRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addColonToTextRule(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeColonAndTextRule()
		})
	}
}

func Test_mergeToRootNode(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeToRootNode()
		})
	}
}

func Test_parser_getPreviousLexeme(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		p    *parser
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.getPreviousLexeme(tt.args.n); got != tt.want {
				t.Errorf("parser.getPreviousLexeme() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newParser(t *testing.T) {
	type args struct {
		lexemes []lexer.Lexeme
	}
	tests := []struct {
		name string
		args args
		want *parser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newParser(tt.args.lexemes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_getCurrentLexeme(t *testing.T) {
	tests := []struct {
		name string
		p    *parser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.getCurrentLexeme()
		})
	}
}
