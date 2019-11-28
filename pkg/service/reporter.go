package service

import (
	"log"
	"time"

	"github.com/dnguy078/healthcheck/pkg/models"
)

// Reporter schedules the healthcheck based on checkFrequency
type Reporter struct {
	tickRate   time.Duration
	quit       chan bool
	results    chan *models.HealthCheck
	dispatcher *Dispatcher
	storage    hcStorage
}

type hcStorage interface {
	List() models.HealthChecks
	Create(input *models.HealthCheck) error
}

// NewReporter returns a reporter
func NewReporter(frequencyRate time.Duration, db hcStorage) (*Reporter, error) {
	incoming := make(chan *models.HealthCheck)
	results := make(chan *models.HealthCheck)
	d, err := NewDispatcher(incoming, results)
	if err != nil {
		return nil, err
	}
	d.Run()

	return &Reporter{
		tickRate:   frequencyRate,
		quit:       make(chan bool),
		dispatcher: d,
		results:    results,
		storage:    db,
	}, nil
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
					r.dispatcher.Process(hc)
				}
			case res := <-r.results:
				go r.storage.Create(res)
			case <-r.quit:
				ticker.Stop()
				r.dispatcher.Stop()
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
