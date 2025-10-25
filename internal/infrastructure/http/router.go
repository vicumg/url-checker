package http

import (
	"net/http"
	"urlChecker/internal/interface/api"
)

func NewRouter(handler *api.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/monitors", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateMonitor(w, r)
		case http.MethodGet:
			handler.GetAllMonitors(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("GET /monitors/{id}", handler.GetMonitor)
	mux.HandleFunc("PUT /monitors/{id}", handler.UpdateMonitor)
	mux.HandleFunc("DELETE /monitors/{id}", handler.DeleteMonitor)
	mux.HandleFunc("POST /monitors/{id}/resume", handler.ResumeMonitor)
	mux.HandleFunc("POST /monitors/{id}/pause", handler.PauseMonitor)
	return mux
}
