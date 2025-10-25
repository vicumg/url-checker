package service

import (
	"time"
	"urlChecker/internal/domain/monitor"
)

type MonitorService struct {
	repo monitor.Repository
}

func NewMonitorService(repo monitor.Repository) *MonitorService {
	return &MonitorService{repo: repo}
}

func (s *MonitorService) CreateMonitor(url string, intervalMinutes int) (*monitor.URLMonitor, error) {
	m := monitor.NewURLMonitor(url, time.Duration(intervalMinutes)*time.Minute)
	err := s.repo.Save(m)
	return m, err
}

func (s *MonitorService) GetMonitor(id string) (*monitor.URLMonitor, error) {
	return s.repo.FindByID(id)
}

func (s *MonitorService) GetAllMonitors() ([]*monitor.URLMonitor, error) {
	return s.repo.FindAll()
}

func (s *MonitorService) UpdateMonitor(id, url string, intervalMinutes int) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	m.Update(url, time.Duration(intervalMinutes)*time.Minute)
	return s.repo.Update(m)
}

func (s *MonitorService) DeleteMonitor(id string) error {
	return s.repo.Delete(id)
}

func (s *MonitorService) PauseMonitor(id string) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	m.Pause()
	return s.repo.Update(m)
}

func (s *MonitorService) ResumeMonitor(id string) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	m.Resume()
	return s.repo.Update(m)
}
