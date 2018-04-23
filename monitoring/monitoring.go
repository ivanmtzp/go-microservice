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

type ReadyStatus struct {
	Database string `json:"database"`
}

func healthinessHandler() func(http.ResponseWriter, *http.Request) {
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

func readinessHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		bytes, err := json.MarshalIndent(&ReadyStatus{Database: "ok"}, "", "\t")
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

type StatusServer struct {
	address string
}

func NewStatusServer(address string) *StatusServer {
	return &StatusServer{address: address}
}

func (s* StatusServer) Address() string {
	return s.address
}

func (s *StatusServer) Run() {
	http.HandleFunc("/healthy", healthinessHandler())
	http.HandleFunc("/ready", readinessHandler())
	http.HandleFunc("/metrics", metricsHandler())
	http.ListenAndServe(s.address, nil)
}