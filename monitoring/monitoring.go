package monitoring

import (
	"net/http"
	"encoding/json"
)

type HealthStatus struct {
	database string `json:"database"`
}

func handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		bytes, err := json.MarshalIndent(&HealthStatus{database: "ok"}, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(bytes)
	}
}

type Server struct {
	enabled bool
}

func (sm *Server) Run() {
	http.HandleFunc("/health", handler())
	http.ListenAndServe("localhost:8080", nil)
}