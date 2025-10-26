package repository

import (
	"database/sql"
	"errors"
	"time"
	"urlChecker/internal/domain/monitor"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &SQLiteRepository{db: db}
	if err := repo.createTable(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SQLiteRepository) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS monitors (
		id TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		interval_seconds INTEGER NOT NULL,
		is_active INTEGER NOT NULL,
		last_checked INTEGER,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	)`

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Save(m *monitor.URLMonitor) error {
	query := `
	INSERT INTO monitors (id, url, interval_seconds, is_active, last_checked, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	var lastChecked *int64
	if m.LastChecked != nil {
		ts := m.LastChecked.Unix()
		lastChecked = &ts
	}

	_, err := r.db.Exec(query,
		m.ID,
		m.URL,
		int64(m.Interval.Seconds()),
		boolToInt(m.IsActive),
		lastChecked,
		m.CreatedAt.Unix(),
		m.UpdatedAt.Unix(),
	)

	return err
}

func (r *SQLiteRepository) FindByID(id string) (*monitor.URLMonitor, error) {
	query := `
	SELECT id, url, interval_seconds, is_active, last_checked, created_at, updated_at
	FROM monitors WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var m monitor.URLMonitor
	var intervalSeconds int64
	var isActive int
	var lastChecked *int64
	var createdAt, updatedAt int64

	err := row.Scan(&m.ID, &m.URL, &intervalSeconds, &isActive, &lastChecked, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("monitor not found")
	}
	if err != nil {
		return nil, err
	}

	m.Interval = time.Duration(intervalSeconds) * time.Second
	m.IsActive = intToBool(isActive)
	if lastChecked != nil {
		t := time.Unix(*lastChecked, 0)
		m.LastChecked = &t
	}
	m.CreatedAt = time.Unix(createdAt, 0)
	m.UpdatedAt = time.Unix(updatedAt, 0)

	return &m, nil
}

func (r *SQLiteRepository) FindAll() ([]*monitor.URLMonitor, error) {
	query := `
	SELECT id, url, interval_seconds, is_active, last_checked, created_at, updated_at
	FROM monitors`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []*monitor.URLMonitor

	for rows.Next() {
		var m monitor.URLMonitor
		var intervalSeconds int64
		var isActive int
		var lastChecked *int64
		var createdAt, updatedAt int64

		err := rows.Scan(&m.ID, &m.URL, &intervalSeconds, &isActive, &lastChecked, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		m.Interval = time.Duration(intervalSeconds) * time.Second
		m.IsActive = intToBool(isActive)
		if lastChecked != nil {
			t := time.Unix(*lastChecked, 0)
			m.LastChecked = &t
		}
		m.CreatedAt = time.Unix(createdAt, 0)
		m.UpdatedAt = time.Unix(updatedAt, 0)

		monitors = append(monitors, &m)
	}

	return monitors, rows.Err()
}

func (r *SQLiteRepository) Delete(id string) error {
	query := `DELETE FROM monitors WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteRepository) Update(m *monitor.URLMonitor) error {
	query := `
	UPDATE monitors
	SET url = ?, interval_seconds = ?, is_active = ?, last_checked = ?, updated_at = ?
	WHERE id = ?`

	var lastChecked *int64
	if m.LastChecked != nil {
		ts := m.LastChecked.Unix()
		lastChecked = &ts
	}

	_, err := r.db.Exec(query,
		m.URL,
		int64(m.Interval.Seconds()),
		boolToInt(m.IsActive),
		lastChecked,
		m.UpdatedAt.Unix(),
		m.ID,
	)

	return err
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func intToBool(i int) bool {
	return i != 0
}
