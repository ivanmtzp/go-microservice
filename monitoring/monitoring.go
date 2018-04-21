package monitoring

import (
	"net/http"
	"encoding/json"
	"fmt"
	"bytes"
)

type HealthStatus struct {
	Database string `json:"database"`
}

func healthCheckHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		bytes, err := json.MarshalIndent(&HealthStatus{Database: "ok"}, "", "\t")
		fmt.Print("Esto:", string(bytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(bytes)
	}
}

func metricsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var buffer bytes.Buffer
		WriteJsonMetrics(&buffer)
		w.Write(buffer.Bytes())
	}
}

type Server struct {
	address string
}

func New(address string) *Server {
	return &Server{address: address}
}

func (s* Server) Address() string {
	return s.address
}

func (s *Server) Run() {
	http.HandleFunc("/healthy", healthCheckHandler())
	http.HandleFunc("/metrics", metricsHandler())
	http.ListenAndServe(s.address, nil)
}