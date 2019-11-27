package models

type HealthCheck struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Code     int    `json:"code"`
	Endpoint string `json:"endpoint"`
	Checked  int    `json:"checked"`
	Duration string `json:"duration"`
}

type CreateHealthCheck struct {
}

type AutoGenerated struct {
	ID       string `json:"id"`
	Endpoint string `json:"endpoint"`
}
