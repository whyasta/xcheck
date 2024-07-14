package constant

import (
	"testing"
)

func TestResponseStatus_GetResponseStatus(t *testing.T) {
	tests := []struct {
		name string
		r    ResponseStatus
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetResponseStatus(); got != tt.want {
				t.Errorf("ResponseStatus.GetResponseStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseStatus_GetResponseMessage(t *testing.T) {
	tests := []struct {
		name string
		r    ResponseStatus
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetResponseMessage(); got != tt.want {
				t.Errorf("ResponseStatus.GetResponseMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
