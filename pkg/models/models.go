package models

import (
	"github.com/dnguy078/healthcheck/pkg/utils"
)

type HealthCheck struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Code     int32  `json:"code"`
	Endpoint string `json:"endpoint"`
	Checked  int64  `json:"checked"`
	Duration string `json:"duration"`
	Error    string `json:"error,omitempty"`
}

func NewHealthCheck(endpoint string) (*HealthCheck, error) {
	uuid, err := utils.UUID()
	if err != nil {
		return nil, err
	}

	return &HealthCheck{
		Endpoint: endpoint,
		ID:       uuid,
	}, nil
}

type HealthChecks []*HealthCheck

func (hcs HealthChecks) Len() int {
	return len(hcs)
}
func (hcs HealthChecks) Swap(i, j int) {
	hcs[i], hcs[j] = hcs[j], hcs[i]
}
func (hcs HealthChecks) Less(i, j int) bool {
	return hcs[i].Endpoint < hcs[j].Endpoint
}

type HealthCheckList struct {
	Items HealthChecks `json:"items"`
	Page  int          `json:"page"`
	Total int          `json:"total"`
	Size  int          `json:"size"`
}

type CreateHealthCheckResponse struct {
	ID       string `json:"id"`
	Endpoint string `json:"endpoint"`
}

type CreateHealthCheckRequest struct {
	Endpoint string `json:"endpoint"`
}
