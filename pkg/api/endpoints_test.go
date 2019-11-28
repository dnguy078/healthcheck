package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/dnguy078/healthcheck/pkg/models"
	"github.com/dnguy078/healthcheck/pkg/storage/mocks"
)

func TestHealthCheckHandler_List(t *testing.T) {
	type fields struct {
		db healthCheckStorage
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		expectedStatusCode int
		want               models.HealthCheckList
	}{
		{
			name: "success",
			fields: fields{
				db: &mocks.FakeCollection{
					ListResp: models.HealthChecks{
						&models.HealthCheck{
							Endpoint: "b",
						},
						&models.HealthCheck{
							Endpoint: "a",
						},
					},
				},
			},
			want: models.HealthCheckList{
				Items: models.HealthChecks{
					&models.HealthCheck{
						Endpoint: "a",
					},
					&models.HealthCheck{
						Endpoint: "b",
					},
				},
				Page:  0,
				Size:  10,
				Total: 2,
			},
			expectedStatusCode: http.StatusOK,
			args: args{
				r: httptest.NewRequest("GET", "/api/health/checks?page=0", nil),
			},
		}, {
			name: "page error",
			args: args{
				r: httptest.NewRequest("GET", "/api/health/checks?page=sdfdsf", nil),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hh := &HealthCheckHandler{
				db: tt.fields.db,
			}
			w := httptest.NewRecorder()

			hh.List(w, tt.args.r)
			if tt.expectedStatusCode != w.Code {
				t.Errorf("got statuscode %d expected code %d", w.Code, tt.expectedStatusCode)
				return
			}

			if len(tt.want.Items) != 0 {
				expected, _ := json.Marshal(tt.want)
				if !reflect.DeepEqual(w.Body.Bytes(), expected) {
					t.Errorf("expected to be equal \n%s\n%s", string(w.Body.String()), string(expected))
				}
			}
		})
	}
}

func TestHealthCheckHandler_Create(t *testing.T) {
	type fields struct {
		db healthCheckStorage
	}
	tests := []struct {
		name               string
		fields             fields
		expectedStatusCode int
		payload            string
	}{
		{
			name: "success",
			fields: fields{
				db: &mocks.FakeCollection{},
			},
			payload:            `{"endpoint":  "https://www.blizzard.com/en-us/"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "empty endpoint",
			fields: fields{
				db: &mocks.FakeCollection{},
			},
			payload:            `{"endpoint":  ""}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "marshall err",
			fields: fields{
				db: &mocks.FakeCollection{},
			},
			payload:            `<<<`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid url",
			fields: fields{
				db: &mocks.FakeCollection{},
			},
			payload:            `{"endpoint":  "ww12bliz.it"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "db error",
			fields: fields{
				db: &mocks.FakeCollection{
					CreateErr: errors.New("duplicate"),
				},
			},
			payload:            `{"endpoint":  "https://www.blizzard.com/en-us/"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hh := &HealthCheckHandler{
				db: tt.fields.db,
			}
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/health/checks", strings.NewReader(tt.payload))
			hh.Create(w, req)

			if tt.expectedStatusCode != w.Code {
				t.Errorf("got statuscode %d expected code %d", w.Code, tt.expectedStatusCode)
			}
		})
	}
}

func TestHealthCheckHandler_Execute(t *testing.T) {
	type fields struct {
		db healthCheckStorage
	}
	tests := []struct {
		name               string
		fields             fields
		expectedStatusCode int
		url                string
	}{
		{
			name: "success",
			fields: fields{
				db: &mocks.FakeCollection{
					GetResp: &models.HealthCheck{
						Endpoint: "https://www.google.com",
					},
				},
			},
			url:                "/api/health/checks/C6C5B3DC-6685-7698-3CD5-C3AB7C10B3AC/try?timeout=2s",
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO set up a fake httpserver instead of hitting google
			hh := &HealthCheckHandler{
				db: tt.fields.db,
			}
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", tt.url, nil)
			hh.Execute(w, req)

			if tt.expectedStatusCode != w.Code {
				t.Errorf("got statuscode %d expected code %d", w.Code, tt.expectedStatusCode)
			}
		})
	}
}
