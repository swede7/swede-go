package model

import "testing"

func TestGetVariableTypeByName(t *testing.T) {
	tests := []struct {
		name    string
		want    VariableType
		wantErr bool
	}{
		{name: "string", want: String, wantErr: false},
		{name: "int", want: Int, wantErr: false},
		{name: "incorrect type", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GetVariableTypeByName(test.name)
			if (err != nil) != test.wantErr {
				t.Errorf("GetVariableTypeByName() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("GetVariableTypeByName() got = %v, want %v", got, test.want)
			}
		})
	}
}
