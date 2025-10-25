package monitor

import (
	"testing"
	"time"
)

func TestNewURLMonitor(t *testing.T) {
	url := "https://example.com"
	interval := 5 * time.Minute

	m := NewURLMonitor(url, interval)

	if m.URL != url {
		t.Errorf("expected URL %s, got %s", url, m.URL)
	}
	if m.Interval != interval {
		t.Errorf("expected interval %v, got %v", interval, m.Interval)
	}
	if !m.IsActive {
		t.Error("expected monitor to be active")
	}
	if m.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestURLMonitor_Pause(t *testing.T) {
	m := NewURLMonitor("https://example.com", 5*time.Minute)

	m.Pause()

	if m.IsActive {
		t.Error("expected monitor to be inactive after pause")
	}
}

func TestURLMonitor_Resume(t *testing.T) {
	m := NewURLMonitor("https://example.com", 5*time.Minute)
	m.Pause()

	m.Resume()

	if !m.IsActive {
		t.Error("expected monitor to be active after resume")
	}
}

func TestURLMonitor_Update(t *testing.T) {
	m := NewURLMonitor("https://example.com", 5*time.Minute)
	newURL := "https://newexample.com"
	newInterval := 10 * time.Minute

	m.Update(newURL, newInterval)

	if m.URL != newURL {
		t.Errorf("expected URL %s, got %s", newURL, m.URL)
	}
	if m.Interval != newInterval {
		t.Errorf("expected interval %v, got %v", newInterval, m.Interval)
	}
}
