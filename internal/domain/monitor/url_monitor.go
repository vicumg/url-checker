package monitor

import (
	"time"
)

type URLMonitor struct {
	ID          string
	URL         string
	Interval    time.Duration
	IsActive    bool
	LastChecked *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewURLMonitor(url string, interval time.Duration) *URLMonitor {
	now := time.Now()
	return &URLMonitor{
		ID:        generateID(),
		URL:       url,
		Interval:  interval,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *URLMonitor) Pause() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

func (u *URLMonitor) Resume() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

func (u *URLMonitor) Update(url string, interval time.Duration) {
	u.URL = url
	u.Interval = interval
	u.UpdatedAt = time.Now()
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
