package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/dnguy078/healthcheck/pkg/models"
)

var (
	testLoadFilePath = "./testdata/seed_load_test.json"
	testDumpFilePath = "./testdata/seed_dump_test.json"
)

func TestCollection_Load(t *testing.T) {
	tests := []struct {
		name            string
		numHealthChecks int
	}{
		{
			name:            "success",
			numHealthChecks: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollection(testLoadFilePath)
			if len(c.List()) != 1 {
				t.Error("expected to load 1 healthcheck")
			}
		})
	}
}

func TestCollection_Dump(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(testDumpFilePath)

			c := NewCollection(testDumpFilePath)
			c.Create(&models.HealthCheck{ID: "testID"})
			c.Dump(testDumpFilePath)

			b, err := ioutil.ReadFile(testDumpFilePath)
			if err != nil {
				t.Error("expected to be able to open dumped file")
			}
			if len(b) == 0 {
				t.Error("expected to have written to desk")
			}
		})
	}
}
