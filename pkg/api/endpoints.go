package endpoints

type HealthCheckHandler struct {
	db healthCheckStorage
}

type healthCheckStorage interface {
	GetHealthChecks()
}