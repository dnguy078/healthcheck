package mocks

import "github.com/dnguy078/healthcheck/pkg/models"

type FakeCollection struct {
	ListResp     models.HealthChecks
	GetResp      *models.HealthCheck
	GetErr       error
	CreateErr    error
	CalledDelete bool
}

func (fc *FakeCollection) List() models.HealthChecks {
	return fc.ListResp
}

func (fc *FakeCollection) Get(id string) (*models.HealthCheck, error) {
	return fc.GetResp, fc.GetErr
}

func (fc *FakeCollection) Create(*models.HealthCheck) error {
	return fc.CreateErr
}

func (fc *FakeCollection) Delete(id string) {
	fc.CalledDelete = true
}
