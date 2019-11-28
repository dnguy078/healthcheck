package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/dnguy078/healthcheck/pkg/models"
	"github.com/dnguy078/healthcheck/pkg/service"
	"github.com/dnguy078/healthcheck/pkg/utils"
)

type HealthCheckHandler struct {
	db healthCheckStorage
}

type healthCheckStorage interface {
	List() models.HealthChecks
	Get(id string) (*models.HealthCheck, error)
	Create(*models.HealthCheck) error
	Delete(id string)
}

func (hh *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	switch r.Method {
	case http.MethodGet:
		if utils.ContainsUUID(r.URL.String()) {
			hh.Get(w, r)
			return
		}
		if _, ok := queryParams["page"]; ok {
			hh.List(w, r)
			return
		}
	case http.MethodPost:
		if utils.ContainsUUID(r.URL.String()) {
			hh.Execute(w, r)
			return
		}
		hh.Create(w, r)
		return
	case http.MethodDelete:
		hh.Delete(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// List returns a paginated list of healthchecks
func (hh *HealthCheckHandler) List(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, err := strconv.Atoi(queryParams["page"][0])
	if err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusBadRequest)
		return
	}

	list := hh.db.List()
	sort.Sort(models.HealthChecks(list))
	start, end := paginate(page, 10, len(list))
	paginated := list[start:end]

	res := &models.HealthCheckList{
		Items: paginated,
		Total: len(list),
		Page:  page,
		Size:  10,
	}

	b, err := json.Marshal(res)
	if err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

// Get returns a specific healthcheck
func (hh *HealthCheckHandler) Get(w http.ResponseWriter, r *http.Request) {
	uuid := utils.ExtractUUID(r.URL.String())
	hc, err := hh.db.Get(uuid)
	if err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(hc)
	if err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

// Create creates a new healthcheck
func (hh *HealthCheckHandler) Create(w http.ResponseWriter, r *http.Request) {
	req := &models.CreateHealthCheckRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusBadRequest)
		return
	}

	if req.Endpoint == "" {
		http.Error(w, marshalError("empty healthcheck endpoint"), http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(req.Endpoint); err != nil {
		http.Error(w, marshalError("invalid URL"), http.StatusBadRequest)
		return
	}

	hc, err := models.NewHealthCheck(req.Endpoint)
	if err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusInternalServerError)
		return
	}

	if err := hh.db.Create(hc); err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusBadRequest)
		return
	}

	resp := &models.CreateHealthCheckResponse{
		ID:       hc.ID,
		Endpoint: hc.Endpoint,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

// Delete removes a healthcheck
func (hh *HealthCheckHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uuid := utils.ExtractUUID(r.URL.String())
	hh.db.Delete(uuid)
}

// Execute a healthcheck with a timeouts
func (hh *HealthCheckHandler) Execute(w http.ResponseWriter, r *http.Request) {
	uuid := utils.ExtractUUID(r.URL.String())
	if uuid == "" {
		http.Error(w, marshalError("invalid uuid"), http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	timeout, err := time.ParseDuration(queryParams["timeout"][0])
	if err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusBadRequest)
		return
	}

	hc, err := hh.db.Get(uuid)
	if err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusBadRequest)
		return
	}

	// make a copy and run the healthcheck
	try := &models.HealthCheck{
		ID:       hc.ID,
		Endpoint: hc.Endpoint,
	}

	try = service.Run(try, timeout)

	b, err := json.Marshal(try)
	if err != nil {
		http.Error(w, marshalError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

// marshalError wraps a error
func marshalError(errString string) string {
	type endpointError struct {
		Message string `json:"error"`
	}
	e := &endpointError{Message: errString}
	b, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(b)
}

func paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
	start := pageNum * pageSize

	if start > sliceLength {
		start = sliceLength
	}

	end := start + pageSize
	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}
