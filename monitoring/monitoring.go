package monitoring

import (
	"net/http"
	"encoding/json"
	"fmt"
)

type HealthStatus struct {
	database string `json:"database"`
}

func handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		bytes, err := json.MarshalIndent(&HealthStatus{database: "ok"}, "", "\t")
		fmt.Print("Esto:", string(bytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(bytes)
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
	http.HandleFunc("/healthy", handler())
	http.ListenAndServe(s.address, nil)
}