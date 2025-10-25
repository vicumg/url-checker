package api

import (
	"encoding/json"
	"net/http"
	"urlChecker/internal/application/service"
)

type Handler struct {
	service *service.MonitorService
}

func NewHandler(service *service.MonitorService) *Handler {
	return &Handler{service: service}
}

type CreateMonitorRequest struct {
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}

type UpdateMonitorRequest struct {
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}

func (h *Handler) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	var req CreateMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := h.service.CreateMonitor(req.URL, req.Interval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func (h *Handler) GetAllMonitors(w http.ResponseWriter, r *http.Request) {
	monitors, err := h.service.GetAllMonitors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monitors)
}

func (h *Handler) GetMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := h.service.GetMonitor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func (h *Handler) UpdateMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdateMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.UpdateMonitor(id, req.URL, req.Interval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.service.DeleteMonitor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PauseMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.service.PauseMonitor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ResumeMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.service.ResumeMonitor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
