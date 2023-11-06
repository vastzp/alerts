package server

import (
	"context"
	"github.com/vastzp/alerts/service"
	"log"
	"net/http"
)

// server is a structure for the server layer.
type server struct {
	service service.Service
}

// NewServer creates a new server instance (constructor).
func NewServer(service service.Service) *server {
	return &server{
		service: service,
	}
}

// Run runs the server. I hardcoded port 8088.
func (s *server) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleGetAlerts(w, r)
		case http.MethodPost:
			s.handlePostAlert(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(":8088", mux)
	if err != nil {
		log.Fatalf("failed to start HTTP server: %s", err)
		return
	}
}

func (s *server) Shutdown(_ context.Context) error {
	// todo: just stub for now
	return nil
}
