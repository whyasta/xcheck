package controllers

import (
	"bigmind/xcheck-be/internal/services"
	"reflect"
	"testing"
)

func TestNewController(t *testing.T) {
	type args struct {
		services *services.Service
	}
	tests := []struct {
		name string
		args args
		want *Controller
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewController(tt.args.services); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewController() = %v, want %v", got, tt.want)
			}
		})
	}
}
