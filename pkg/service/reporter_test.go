package service

import (
	"testing"
	"time"

	"github.com/dnguy078/healthcheck/pkg/storage/mocks"
)

func TestNewReporter(t *testing.T) {
	type args struct {
		frequencyRate time.Duration
		db            hcStorage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				frequencyRate: 1 * time.Microsecond,
				db:            &mocks.FakeCollection{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewReporter(tt.args.frequencyRate, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReporter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			r.Report()
			r.Stop()
		})
	}
}
