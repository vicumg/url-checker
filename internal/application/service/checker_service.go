package service

import (
	"context"
	"log"
	"net/http"
	"time"
	"urlChecker/internal/domain/monitor"
)

type Logger interface {
	LogCheck(monitorID, url string, statusCode int, responseTime time.Duration, err error)
}

type CheckerService struct {
	repo   monitor.Repository
	client *http.Client
	logger Logger
}

func NewCheckerService(repo monitor.Repository, logger Logger) *CheckerService {
	return &CheckerService{
		repo:   repo,
		logger: logger,
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
	start := time.Now()
	resp, err := s.client.Get(m.URL)
	responseTime := time.Since(start)

	now := time.Now()
	m.LastChecked = &now

	if err != nil {
		s.logger.LogCheck(m.ID, m.URL, 0, responseTime, err)
	} else {
		s.logger.LogCheck(m.ID, m.URL, resp.StatusCode, responseTime, nil)
		resp.Body.Close()
	}

	s.repo.Update(m)
}
