package service

import (
	"context"
	"log"
	"net/http"
	"time"
	"urlChecker/internal/domain/monitor"
)

type CheckerService struct {
	repo   monitor.Repository
	client *http.Client
}

func NewCheckerService(repo monitor.Repository) *CheckerService {
	return &CheckerService{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *CheckerService) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checkAllMonitors()
		}
	}
}

func (s *CheckerService) checkAllMonitors() {
	monitors, err := s.repo.FindAll()
	if err != nil {
		log.Printf("Error getting monitors: %v", err)
		return
	}

	for _, m := range monitors {
		if !m.IsActive {
			continue
		}

		if m.LastChecked != nil && time.Since(*m.LastChecked) < m.Interval {
			continue
		}

		go s.checkURL(m)
	}
}

func (s *CheckerService) checkURL(m *monitor.URLMonitor) {
	resp, err := s.client.Get(m.URL)
	now := time.Now()
	m.LastChecked = &now

	if err != nil {
		log.Printf("[%s] Error checking %s: %v", m.ID, m.URL, err)
	} else {
		log.Printf("[%s] Checked %s - Status: %d", m.ID, m.URL, resp.StatusCode)
		resp.Body.Close()
	}

	s.repo.Update(m)
}
