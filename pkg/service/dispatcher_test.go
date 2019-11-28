package service

import (
	"reflect"
	"testing"
	"time"

	"github.com/dnguy078/healthcheck/pkg/models"
)

func TestRun(t *testing.T) {
	type args struct {
		hc      *models.HealthCheck
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want *models.HealthCheck
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Run(tt.args.hc, tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
