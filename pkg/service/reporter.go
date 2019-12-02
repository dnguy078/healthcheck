package service

import (
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

// Reporter schedules the healthcheck based on checkFrequency
type Reporter struct {
	tickRate time.Duration
	quit     chan bool
	results  chan *models.HealthCheck
	storage  hcStorage
	jobQueue chan *models.HealthCheck
}

type hcStorage interface {
	List() models.HealthChecks
	Create(input *models.HealthCheck) error
}

// NewReporter returns a reporter
func NewReporter(frequencyRate time.Duration, db hcStorage) (*Reporter, error) {
	results := make(chan *models.HealthCheck)
	r := &Reporter{
		tickRate: frequencyRate,
		quit:     make(chan bool),
		results:  results,
		jobQueue: make(chan *models.HealthCheck),
		storage:  db,
	}

	for i := 0; i < maxWorkers; i++ {
		worker := Worker{
			id:       i,
			jobQueue: r.jobQueue,
			quit:     r.quit,
			results:  r.results,
		}

		go worker.Start()
	}

	return r, nil
}

// Report ticks based upon check frequency and routes healthchecks to be performed to the dispatcher
func (r *Reporter) Report() {
	ticker := time.NewTicker(r.tickRate)
	go func() {
		for {
			select {
			case <-ticker.C:
				list := r.storage.List()

				for _, hc := range list {
					r.jobQueue <- hc
				}
			case res := <-r.results:
				go r.storage.Create(res)
			case <-r.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop the reporter
func (r *Reporter) Stop() {
	log.Print("Stopping healthcheck reporter")
	close(r.quit)
}
