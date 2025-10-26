package repository

import (
	"os"
	"testing"
	"time"
	"urlChecker/internal/domain/monitor"
)

func TestSQLiteRepository_Save(t *testing.T) {
	dbPath := "test_save.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close()

	m := monitor.NewURLMonitor("https://example.com", 5*time.Minute)
	err = repo.Save(m)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSQLiteRepository_FindByID(t *testing.T) {
	dbPath := "test_find.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close()

	m := monitor.NewURLMonitor("https://example.com", 5*time.Minute)
	repo.Save(m)

	found, err := repo.FindByID(m.ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if found.ID != m.ID {
		t.Errorf("expected ID %s, got %s", m.ID, found.ID)
	}
	if found.URL != m.URL {
		t.Errorf("expected URL %s, got %s", m.URL, found.URL)
	}
}

func TestSQLiteRepository_FindAll(t *testing.T) {
	dbPath := "test_findall.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close()

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

func TestSQLiteRepository_Update(t *testing.T) {
	dbPath := "test_update.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close()

	m := monitor.NewURLMonitor("https://example.com", 5*time.Minute)
	repo.Save(m)

	m.Update("https://updated.com", 10*time.Minute)
	err = repo.Update(m)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	updated, _ := repo.FindByID(m.ID)
	if updated.URL != "https://updated.com" {
		t.Errorf("expected URL to be updated to https://updated.com, got %s", updated.URL)
	}
}

func TestSQLiteRepository_Delete(t *testing.T) {
	dbPath := "test_delete.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close()

	m := monitor.NewURLMonitor("https://example.com", 5*time.Minute)
	repo.Save(m)

	err = repo.Delete(m.ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = repo.FindByID(m.ID)
	if err == nil {
		t.Error("expected monitor to be deleted")
	}
}
