package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/dnguy078/healthcheck/pkg/models"
)

// Collection holds a map of healthchecks
type Collection struct {
	// maps are not thread safe in golang; have to add mutex around this map
	sync.RWMutex
	data           map[string]*models.HealthCheck
	registeredURLs map[string]bool
}

// NewCollection returns a new Collection
func NewCollection(filePath string) *Collection {
	c := &Collection{
		data:           map[string]*models.HealthCheck{},
		registeredURLs: make(map[string]bool),
	}
	if err := c.Load(filePath); err != nil {
		log.Printf("unable to load any existing healthchecks from disk, err: %s", err)
	}

	return c
}

// List returns a list of healthchecks
func (c *Collection) List() models.HealthChecks {
	items := make([]*models.HealthCheck, 0)
	c.RLock()
	c.RUnlock()
	for _, c := range c.data {
		items = append(items, c)
	}

	return items
}

// Get returns a specific heallthcheck, errors if healthcheck does not exist
func (c *Collection) Get(id string) (*models.HealthCheck, error) {
	c.RLock()
	defer c.RUnlock()
	hc, ok := c.data[id]
	if !ok {
		return nil, fmt.Errorf("healthcheck %s not found", id)
	}

	return hc, nil
}

// Create adds a healthcheck to the collection
func (c *Collection) Create(input *models.HealthCheck) error {
	c.Lock()
	defer c.Unlock()
	if _, found := c.registeredURLs[input.Endpoint]; !found {
		c.data[input.ID] = input
		c.registeredURLs[input.Endpoint] = true
		return nil
	}
	return fmt.Errorf("endpoint %s already registered", input.Endpoint)
}

// Delete removes a healthcheck from the collection
func (c *Collection) Delete(id string) {
	c.Lock()
	defer c.Unlock()
	if hc, found := c.data[id]; found {
		delete(c.registeredURLs, hc.Endpoint)
		delete(c.data, id)
	}
}

// Dump takes all existing healthchecks and writes them to disk in JSON format
func (c *Collection) Dump(fileName string) error {
	list := c.List()
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}

	f, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	return err
}

// Load reads a file container healthchecks written as JSON and populates the collection with existing healthchecks
func (c *Collection) Load(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		return err
	}

	list := make([]*models.HealthCheck, 0)

	if err := json.Unmarshal(b, &list); err != nil {
		return err
	}

	for _, h := range list {
		c.data[h.ID] = h
	}

	return nil
}
