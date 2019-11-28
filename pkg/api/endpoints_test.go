package api

import (
	"net/http"
	"testing"
)

func TestHealthCheckHandler_List(t *testing.T) {
	type fields struct {
		db healthCheckStorage
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hh := &HealthCheckHandler{
				db: tt.fields.db,
			}
			hh.List(tt.args.w, tt.args.r)
		})
	}
}
