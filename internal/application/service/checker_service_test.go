package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"urlChecker/internal/domain/monitor"
	"urlChecker/internal/infrastructure/repository"
)

// Mock Logger для тестов
type MockLogger struct {
	logs []LogEntry
}

type LogEntry struct {
	MonitorID    string
	URL          string
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

func (m *MockLogger) LogCheck(monitorID, url string, statusCode int, responseTime time.Duration, err error) {
	m.logs = append(m.logs, LogEntry{
		MonitorID:    monitorID,
		URL:          url,
		StatusCode:   statusCode,
		ResponseTime: responseTime,
		Error:        err,
	})
}

func TestCheckerService_CheckURL_Success(t *testing.T) {
	// Создаем тестовый HTTP сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	repo := repository.NewMemoryRepository()
	mockLogger := &MockLogger{}
	checker := NewCheckerService(repo, mockLogger)

	m := monitor.NewURLMonitor(server.URL, 1*time.Minute)
	repo.Save(m)

	// Проверяем URL
	checker.checkURL(m)

	// Проверяем что логирование произошло
	if len(mockLogger.logs) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(mockLogger.logs))
	}

	log := mockLogger.logs[0]
	if log.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", log.StatusCode)
	}
	if log.Error != nil {
		t.Errorf("expected no error, got %v", log.Error)
	}
}

func TestCheckerService_CheckURL_Error(t *testing.T) {
	repo := repository.NewMemoryRepository()
	mockLogger := &MockLogger{}
	checker := NewCheckerService(repo, mockLogger)

	// Невалидный URL
	m := monitor.NewURLMonitor("http://invalid-url-that-does-not-exist-12345.com", 1*time.Minute)
	repo.Save(m)

	checker.checkURL(m)

	if len(mockLogger.logs) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(mockLogger.logs))
	}

	log := mockLogger.logs[0]
	if log.Error == nil {
		t.Error("expected error, got nil")
	}
}

func TestCheckerService_CheckAllMonitors_SkipsInactive(t *testing.T) {
	repo := repository.NewMemoryRepository()
	mockLogger := &MockLogger{}
	checker := NewCheckerService(repo, mockLogger)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	m := monitor.NewURLMonitor(server.URL, 1*time.Minute)
	m.Pause()
	repo.Save(m)

	checker.checkAllMonitors()

	// Ждем немного для завершения goroutines
	time.Sleep(100 * time.Millisecond)

	// Неактивный монитор не должен проверяться
	if len(mockLogger.logs) != 0 {
		t.Errorf("expected 0 log entries for inactive monitor, got %d", len(mockLogger.logs))
	}
}

func TestCheckerService_CheckAllMonitors_RespectsInterval(t *testing.T) {
	repo := repository.NewMemoryRepository()
	mockLogger := &MockLogger{}
	checker := NewCheckerService(repo, mockLogger)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	m := monitor.NewURLMonitor(server.URL, 10*time.Minute)
	now := time.Now().Add(-5 * time.Minute) // Проверен 5 минут назад
	m.LastChecked = &now
	repo.Save(m)

	checker.checkAllMonitors()

	time.Sleep(100 * time.Millisecond)

	// Интервал не прошел - не должно быть проверки
	if len(mockLogger.logs) != 0 {
		t.Errorf("expected 0 log entries (interval not passed), got %d", len(mockLogger.logs))
	}
}

func TestCheckerService_Start_StopsOnContextCancel(t *testing.T) {
	repo := repository.NewMemoryRepository()
	mockLogger := &MockLogger{}
	checker := NewCheckerService(repo, mockLogger)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan bool)
	go func() {
		checker.Start(ctx)
		done <- true
	}()

	// Отменяем контекст через 100ms
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Ждем завершения Start()
	select {
	case <-done:
		// OK, Start() завершился
	case <-time.After(1 * time.Second):
		t.Error("Start() did not stop after context cancel")
	}
}

func TestCheckerService_CheckAllMonitors_MultipleMonitors(t *testing.T) {
	repo := repository.NewMemoryRepository()
	mockLogger := &MockLogger{}
	checker := NewCheckerService(repo, mockLogger)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Создаем несколько мониторов
	m1 := monitor.NewURLMonitor(server.URL+"/1", 1*time.Minute)
	m2 := monitor.NewURLMonitor(server.URL+"/2", 1*time.Minute)
	m3 := monitor.NewURLMonitor(server.URL+"/3", 1*time.Minute)

	repo.Save(m1)
	repo.Save(m2)
	repo.Save(m3)

	checker.checkAllMonitors()

	time.Sleep(200 * time.Millisecond)

	// Все 3 монитора должны быть проверены
	if len(mockLogger.logs) != 3 {
		t.Errorf("expected 3 log entries, got %d", len(mockLogger.logs))
	}
}
