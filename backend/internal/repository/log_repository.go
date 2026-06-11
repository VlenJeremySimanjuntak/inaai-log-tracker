package repository

import (
	"context"
	"database/sql"
	"fmt"

	"backend-tracker/internal/models"
)

type LogRepository interface {
	FetchAll(ctx context.Context) ([]models.IncidentLog, error)
	Create(ctx context.Context, log *models.IncidentLog) error
	GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int) (*models.IncidentLog, error)
	UpdateStatus(ctx context.Context, tx *sql.Tx, id int, status string) error
	SaveAISummary(ctx context.Context, summary *models.AISummary) error
	GetLatestSummary(ctx context.Context) (*models.AISummary, error)
	BeginTx(ctx context.Context) (*sql.Tx, error) // untuk transaksi di usecase
}

type mysqlLogRepository struct {
	db *sql.DB
}

func NewMysqlLogRepository(db *sql.DB) LogRepository {
	return &mysqlLogRepository{db: db}
}

func (r *mysqlLogRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *mysqlLogRepository) FetchAll(ctx context.Context) ([]models.IncidentLog, error) {
	query := `SELECT id, user_id, category_id, title, description, status, created_at, updated_at FROM incident_logs ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query fetch all: %w", err)
	}
	defer rows.Close()

	var logs []models.IncidentLog
	for rows.Next() {
		var l models.IncidentLog
		err := rows.Scan(&l.ID, &l.UserID, &l.CategoryID, &l.Title, &l.Description, &l.Status, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *mysqlLogRepository) Create(ctx context.Context, log *models.IncidentLog) error {
	query := `INSERT INTO incident_logs (user_id, category_id, title, description, status) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, log.UserID, log.CategoryID, log.Title, log.Description, log.Status)
	if err != nil {
		return fmt.Errorf("insert incident log: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}
	log.ID = int(id)
	return nil
}

func (r *mysqlLogRepository) GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int) (*models.IncidentLog, error) {
	query := `SELECT id, user_id, category_id, title, description, status, created_at, updated_at FROM incident_logs WHERE id = ? FOR UPDATE`
	var l models.IncidentLog
	err := tx.QueryRowContext(ctx, query, id).Scan(&l.ID, &l.UserID, &l.CategoryID, &l.Title, &l.Description, &l.Status, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("log id %d tidak ditemukan", id)
		}
		return nil, fmt.Errorf("select for update: %w", err)
	}
	return &l, nil
}

func (r *mysqlLogRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, id int, status string) error {
	query := `UPDATE incident_logs SET status = ? WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	return nil
}

func (r *mysqlLogRepository) SaveAISummary(ctx context.Context, summary *models.AISummary) error {
	query := `INSERT INTO ai_summaries (summary_text, log_ids_analyzed) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, summary.SummaryText, summary.LogIDsAnalyzed)
	if err != nil {
		return fmt.Errorf("save AI summary: %w", err)
	}
	return nil
}

func (r *mysqlLogRepository) GetLatestSummary(ctx context.Context) (*models.AISummary, error) {
	query := `SELECT id, summary_text, log_ids_analyzed, created_at FROM ai_summaries ORDER BY created_at DESC LIMIT 1`
	var s models.AISummary
	err := r.db.QueryRowContext(ctx, query).Scan(&s.ID, &s.SummaryText, &s.LogIDsAnalyzed, &s.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get latest summary: %w", err)
	}
	return &s, nil
}	