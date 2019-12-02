package service

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnguy078/healthcheck/pkg/models"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int32
		mockServer     bool
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
			mockServer:     true,
		},
		{
			name:           "url is down",
			expectedStatus: 0,
			mockServer:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := startServer(t)

			if tt.mockServer {
				s.Start()
				defer s.Close()
			}

			hc := &models.HealthCheck{
				Endpoint: s.URL,
			}

			got := Run(hc, 1*time.Second)
			if got == nil {
				t.Error("expected to return a healthcheck")
				return
			}
			if got.Code != tt.expectedStatus {
				t.Error("expected to have status OK")
			}
		})
	}
}

func startServer(t *testing.T) *httptest.Server {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	}))
	l, err := net.Listen("tcp", "localhost:1111")
	if err != nil {
		t.Fatal(err)
	}

	ts.Listener = l
	return ts
}
