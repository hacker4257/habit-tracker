package repository

import (
	"database/sql"
	"fmt"
	"time"

	"habit-tracker/internal/config"
	"habit-tracker/internal/model"
	"habit-tracker/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type RecordRepository interface {
	Create(record *model.Record) error
	GetByID(id int64) (*model.Record, error)
	GetAll() ([]model.Record, error)
	Update(record *model.Record) error
	Delete(id int64) error
	GetStats() (*model.Stats, error)
}

type recordRepository struct {
	db *sql.DB
}

func NewRecordRepository(cfg *config.DatabaseConfig) (RecordRepository, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &recordRepository{db: db}
	if err := repo.migrate(cfg.Driver); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	logger.Info("Database connected: %s", cfg.Driver)
	return repo, nil
}

func (r *recordRepository) migrate(driver string) error {
	var schema string
	if driver == "mysql" {
		schema = `
		CREATE TABLE IF NOT EXISTS records (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			date VARCHAR(10) NOT NULL,
			content VARCHAR(255) NOT NULL,
			duration INT NOT NULL,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_date (date)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	} else {
		schema = `
		CREATE TABLE IF NOT EXISTS records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			content TEXT NOT NULL,
			duration INTEGER NOT NULL,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_date ON records(date);`
	}

	_, err := r.db.Exec(schema)
	return err
}

func (r *recordRepository) Create(record *model.Record) error {
	now := time.Now()
	result, err := r.db.Exec(
		`INSERT INTO records (date, content, duration, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		record.Date, record.Content, record.Duration, record.Notes, now, now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	record.ID = id
	record.CreatedAt = now
	record.UpdatedAt = now
	return nil
}

func (r *recordRepository) GetByID(id int64) (*model.Record, error) {
	record := &model.Record{}
	err := r.db.QueryRow(
		`SELECT id, date, content, duration, notes, created_at, updated_at FROM records WHERE id = ?`,
		id,
	).Scan(&record.ID, &record.Date, &record.Content, &record.Duration, &record.Notes, &record.CreatedAt, &record.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *recordRepository) GetAll() ([]model.Record, error) {
	rows, err := r.db.Query(
		`SELECT id, date, content, duration, notes, created_at, updated_at FROM records ORDER BY date DESC, id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.Record
	for rows.Next() {
		var record model.Record
		if err := rows.Scan(&record.ID, &record.Date, &record.Content, &record.Duration, &record.Notes, &record.CreatedAt, &record.UpdatedAt); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, rows.Err()
}

func (r *recordRepository) Update(record *model.Record) error {
	record.UpdatedAt = time.Now()
	result, err := r.db.Exec(
		`UPDATE records SET date = ?, content = ?, duration = ?, notes = ?, updated_at = ? WHERE id = ?`,
		record.Date, record.Content, record.Duration, record.Notes, record.UpdatedAt, record.ID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *recordRepository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM records WHERE id = ?`, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *recordRepository) GetStats() (*model.Stats, error) {
	stats := &model.Stats{}

	// Total records and duration
	err := r.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(duration), 0) FROM records`).
		Scan(&stats.TotalRecords, &stats.TotalDuration)
	if err != nil {
		return nil, err
	}

	// This week
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	weekStartStr := weekStart.Format("2006-01-02")
	err = r.db.QueryRow(`SELECT COUNT(*) FROM records WHERE date >= ?`, weekStartStr).
		Scan(&stats.ThisWeek)
	if err != nil {
		return nil, err
	}

	// This month
	monthStartStr := now.Format("2006-01") + "-01"
	err = r.db.QueryRow(`SELECT COUNT(*) FROM records WHERE date >= ?`, monthStartStr).
		Scan(&stats.ThisMonth)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
