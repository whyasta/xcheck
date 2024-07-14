package server

import (
	"bigmind/xcheck-be/internal/services"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewRouter(t *testing.T) {
	type args struct {
		services *services.Service
	}
	tests := []struct {
		name string
		args args
		want *gin.Engine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRouter(tt.args.services); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}
