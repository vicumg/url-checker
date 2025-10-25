package repository

import (
	"testing"
	"time"
	"urlChecker/internal/domain/monitor"
)

func TestMemoryRepository_Save(t *testing.T) {
	repo := NewMemoryRepository()
	m := monitor.NewURLMonitor("https://example.com", 5*time.Minute)

	err := repo.Save(m)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestMemoryRepository_FindByID(t *testing.T) {
	repo := NewMemoryRepository()
	m := monitor.NewURLMonitor("https://example.com", 5*time.Minute)
	repo.Save(m)

	found, err := repo.FindByID(m.ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if found.ID != m.ID {
		t.Errorf("expected ID %s, got %s", m.ID, found.ID)
	}
}

func TestMemoryRepository_FindByID_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	_, err := repo.FindByID("nonexistent")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestMemoryRepository_FindAll(t *testing.T) {
	repo := NewMemoryRepository()
	m1 := monitor.NewURLMonitor("https://example1.com", 5*time.Minute)
	m2 := monitor.NewURLMonitor("https://example2.com", 10*time.Minute)
	repo.Save(m1)
	repo.Save(m2)

	all, err := repo.FindAll()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 monitors, got %d", len(all))
	}
}

func TestMemoryRepository_Delete(t *testing.T) {
	repo := NewMemoryRepository()
	m := monitor.NewURLMonitor("https://example.com", 5*time.Minute)
	repo.Save(m)

	err := repo.Delete(m.ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = repo.FindByID(m.ID)
	if err == nil {
		t.Error("expected monitor to be deleted")
	}
}

func TestMemoryRepository_Update(t *testing.T) {
	repo := NewMemoryRepository()
	m := monitor.NewURLMonitor("https://example.com", 5*time.Minute)
	repo.Save(m)

	m.Update("https://updated.com", 10*time.Minute)
	err := repo.Update(m)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	updated, _ := repo.FindByID(m.ID)
	if updated.URL != "https://updated.com" {
		t.Errorf("expected URL to be updated to https://updated.com, got %s", updated.URL)
	}
}
