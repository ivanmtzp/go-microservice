package monitoring

import (
	"net/http"
)


type StatusServer struct {
	enabled bool
	address string
	healthChecks HealthChecks
}

func NewStatusServer() *StatusServer {
	return &StatusServer{healthChecks: make(HealthChecks)}
}

func (s* StatusServer) Enable(address string) {
	s.enabled = true
	s.address = address
}

func (s* StatusServer) Enabled() bool {
	return s.enabled
}
func (s* StatusServer) Address() string {
	return s.address
}

func (s* StatusServer) RegisterHealthCheck(name string, healthChecker HealthChecker){
	s.healthChecks[name] = healthChecker
}

func (s *StatusServer) Run() {
	if s.enabled {
		http.HandleFunc("/healthy", healthinessHandler(s.healthChecks))
		http.HandleFunc("/metrics", metricsHandler())
		http.ListenAndServe(s.address, nil)
	}
}