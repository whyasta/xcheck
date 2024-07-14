package services

import (
	"bigmind/xcheck-be/internal/repositories"
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {
	type args struct {
		repositories *repositories.Repository
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.args.repositories); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}
