package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dnguy078/healthcheck/pkg/models"
)

var (
	maxWorkers         = 10
	defaultClient      = http.DefaultClient
	defaultHTTPTimeout = 1 * time.Second
)

// Dispatcher spins up a group of goroutines to process healthchecks
type Dispatcher struct {
	jobs   chan *models.HealthCheck
	result chan *models.HealthCheck
	quit   chan bool
}

// NewDispatcher returns a dispatcher
func NewDispatcher(incoming chan *models.HealthCheck, result chan *models.HealthCheck) (*Dispatcher, error) {
	return &Dispatcher{
		jobs:   incoming,
		result: result,
		quit:   make(chan bool),
	}, nil
}

// Process routes message to pool of workers to perform healthchecks
func (d *Dispatcher) Process(hc *models.HealthCheck) {
	d.jobs <- hc
}

// Stop stops pool of workers
func (d *Dispatcher) Stop() {
	close(d.quit)
}

// Run spins up workers
func (d *Dispatcher) Run() {
	for i := 0; i < maxWorkers; i++ {
		worker := Worker{
			id:       i,
			jobQueue: d.jobs,
			quit:     d.quit,
			results:  d.result,
		}

		go worker.Start()
	}
}

// Worker is a goroutine that performs healthchecks
type Worker struct {
	id       int
	jobQueue chan *models.HealthCheck
	results  chan *models.HealthCheck
	quit     chan bool
}

// Start method starts the run loop listening for a job to come in and also listening for a quit signal to stop
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
