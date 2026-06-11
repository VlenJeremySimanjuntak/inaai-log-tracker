package repository

import (
	"context"
	"database/sql"
	"backend-tracker/internal/models"
)

type LogRepository interface {
	FetchAll(ctx context.Context) ([]models.IncidentLog, error)
	Create(ctx context.Context, log *models.IncidentLog) error
	GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int) (*models.IncidentLog, error)
	UpdateStatus(ctx context.Context, tx *sql.Tx, id int, status string) error
	SaveAISummary(ctx context.Context, summary *models.AISummary) error
	GetLatestSummary(ctx context.Context) (*models.AISummary, error)
	GetDB() *sql.DB // Method baru untuk mengakses koneksi database
}

// UBAH: dari mysqlLogRepository menjadi MysqlLogRepository (huruf besar diawal)
type MysqlLogRepository struct {
	Conn *sql.DB
}

// UBAH: return type mengacu ke struct yang sudah diexport
func NewMysqlLogRepository(conn *sql.DB) LogRepository {
	return &MysqlLogRepository{Conn: conn}
}

// Method baru untuk mengakses DB
func (m *MysqlLogRepository) GetDB() *sql.DB {
	return m.Conn
}

func (m *MysqlLogRepository) FetchAll(ctx context.Context) ([]models.IncidentLog, error) {
	query := `SELECT id, user_id, category_id, title, description, status, created_at, updated_at FROM incident_logs ORDER BY created_at DESC`
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.IncidentLog
	for rows.Next() {
		var log models.IncidentLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.CategoryID, &log.Title, &log.Description, &log.Status, &log.CreatedAt, &log.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, log)
	}
	return list, nil
}

func (m *MysqlLogRepository) Create(ctx context.Context, log *models.IncidentLog) error {
	query := `INSERT INTO incident_logs (user_id, category_id, title, description, status) VALUES (?, ?, ?, ?, ?)`
	res, err := m.Conn.ExecContext(ctx, query, log.UserID, log.CategoryID, log.Title, log.Description, log.Status)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	log.ID = int(id)
	return nil
}

// Implementasi PESSIMISTIC LOCKING: Mengunci baris data spesifik agar tidak diubah proses lain selama transaksi berjalan
func (m *MysqlLogRepository) GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int) (*models.IncidentLog, error) {
	query := `SELECT id, user_id, category_id, title, description, status, created_at, updated_at FROM incident_logs WHERE id = ? FOR UPDATE`
	var log models.IncidentLog
	err := tx.QueryRowContext(ctx, query, id).Scan(&log.ID, &log.UserID, &log.CategoryID, &log.Title, &log.Description, &log.Status, &log.CreatedAt, &log.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (m *MysqlLogRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, id int, status string) error {
	query := `UPDATE incident_logs SET status = ? WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, status, id)
	return err
}

func (m *MysqlLogRepository) SaveAISummary(ctx context.Context, summary *models.AISummary) error {
	query := `INSERT INTO ai_summaries (summary_text, log_ids_analyzed) VALUES (?, ?)`
	_, err := m.Conn.ExecContext(ctx, query, summary.SummaryText, summary.LogIDsAnalyzed)
	return err
}

func (m *MysqlLogRepository) GetLatestSummary(ctx context.Context) (*models.AISummary, error) {
	query := `SELECT id, summary_text, log_ids_analyzed, created_at FROM ai_summaries ORDER BY created_at DESC LIMIT 1`
	var summary models.AISummary
	err := m.Conn.QueryRowContext(ctx, query).Scan(&summary.ID, &summary.SummaryText, &summary.LogIDsAnalyzed, &summary.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &summary, nil
}