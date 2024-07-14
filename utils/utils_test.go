package utils

import "testing"

func TestInitLogger(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitLogger(tt.args.env)
		})
	}
}
