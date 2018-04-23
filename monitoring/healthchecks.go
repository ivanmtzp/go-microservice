package monitoring

import (
	"net/http"
	"encoding/json"
)

type HealthChecker interface {
	HealthCheck() error
}

type HealthChecks map[string]HealthChecker

type healthStatus struct {
	Healthy bool `json:"healthy"`
	HealthChecksResults map[string]string `json:"healthchecks"`
}



func healthinessHandler(healthChecks HealthChecks) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hs := &healthStatus{Healthy: true, HealthChecksResults: make(map[string]string)}
		for name, healthCheck := range healthChecks {
			err := healthCheck.HealthCheck()
			if err != nil {
				hs.Healthy = false
				hs.HealthChecksResults[name] = err.Error()
			} else {
				hs.HealthChecksResults[name] = "ok"
			}
		}
		bytes, err := json.MarshalIndent(hs, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(bytes)
	}
}
