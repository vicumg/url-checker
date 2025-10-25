package service

import (
	"testing"
	"time"
	"urlChecker/internal/infrastructure/repository"
)

func TestMonitorService_CreateMonitor(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewMonitorService(repo)

	m, err := service.CreateMonitor("https://example.com", 5)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if m.URL != "https://example.com" {
		t.Errorf("expected URL https://example.com, got %s", m.URL)
	}
	if m.Interval != 5*time.Minute {
		t.Errorf("expected interval 5m, got %v", m.Interval)
	}
}

func TestMonitorService_GetMonitor(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewMonitorService(repo)
	m, _ := service.CreateMonitor("https://example.com", 5)

	found, err := service.GetMonitor(m.ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if found.ID != m.ID {
		t.Errorf("expected ID %s, got %s", m.ID, found.ID)
	}
}

func TestMonitorService_GetAllMonitors(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewMonitorService(repo)
	service.CreateMonitor("https://example1.com", 5)
	service.CreateMonitor("https://example2.com", 10)

	all, err := service.GetAllMonitors()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 monitors, got %d", len(all))
	}
}

func TestMonitorService_UpdateMonitor(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewMonitorService(repo)
	m, _ := service.CreateMonitor("https://example.com", 5)

	err := service.UpdateMonitor(m.ID, "https://updated.com", 10)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	updated, _ := service.GetMonitor(m.ID)
	if updated.URL != "https://updated.com" {
		t.Errorf("expected URL https://updated.com, got %s", updated.URL)
	}
}

func TestMonitorService_DeleteMonitor(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewMonitorService(repo)
	m, _ := service.CreateMonitor("https://example.com", 5)

	err := service.DeleteMonitor(m.ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = service.GetMonitor(m.ID)
	if err == nil {
		t.Error("expected monitor to be deleted")
	}
}

func TestMonitorService_PauseMonitor(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewMonitorService(repo)
	m, _ := service.CreateMonitor("https://example.com", 5)

	err := service.PauseMonitor(m.ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	paused, _ := service.GetMonitor(m.ID)
	if paused.IsActive {
		t.Error("expected monitor to be paused")
	}
}

func TestMonitorService_ResumeMonitor(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewMonitorService(repo)
	m, _ := service.CreateMonitor("https://example.com", 5)
	service.PauseMonitor(m.ID)

	err := service.ResumeMonitor(m.ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	resumed, _ := service.GetMonitor(m.ID)
	if !resumed.IsActive {
		t.Error("expected monitor to be active")
	}
}
