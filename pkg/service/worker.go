package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dnguy078/healthcheck/pkg/models"
)

// Worker is a goroutine that performs healthchecks
type Worker struct {
	id       int
	jobQueue chan *models.HealthCheck
	results  chan *models.HealthCheck
	quit     chan bool
}

// Start method listens for incoming work and runs healthchecks
func (w Worker) Start() {
	for {
		select {
		case job, ok := <-w.jobQueue:
			if !ok {
				return
			}
			w.results <- Run(job, defaultHTTPTimeout)
		case <-w.quit:
			return
		}
	}
}

// Run performs healthchecks
func Run(hc *models.HealthCheck, timeout time.Duration) *models.HealthCheck {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	t := time.Now()
	defer timeRequest(t, hc)

	request, err := http.NewRequest("GET", hc.Endpoint, nil)
	if err != nil {
		handleErr(hc, err)
		return hc
	}
	request = request.WithContext(ctx)

	resp, err := defaultClient.Do(request)
	if err != nil {
		handleErr(hc, err)
		return hc
	}

	hc.Code = int32(resp.StatusCode)
	hc.Status = resp.Status
	hc.Checked = t.Unix()
	hc.Error = ""
	return hc
}

func timeRequest(t time.Time, hc *models.HealthCheck) {
	hc.Duration = time.Since(t).String()
}

func handleErr(hc *models.HealthCheck, err error) {
	hc.Status = "Error"
	hc.Error = err.Error()
	hc.Code = 0
	mapErr := map[string]string{
		"error":   err.Error(),
		"message": "error processing healthcheck",
	}
	log.Print(mapErr)
}
