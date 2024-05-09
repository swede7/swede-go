package highlight

import (
	"reflect"
	"testing"
)

func TestHighlight(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Highlight(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Highlight() = %v, want %v", got, tt.want)
			}
		})
	}
}
