package usecase

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"backend-tracker/internal/models"
	"backend-tracker/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Mock manual untuk testing tanpa database (khusus error di BeginTx)
// ----------------------------------------------------------------------------
type manualMockRepo struct {
	beginTxFunc func(ctx context.Context) (*sql.Tx, error)
}

func (m *manualMockRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if m.beginTxFunc != nil {
		return m.beginTxFunc(ctx)
	}
	return nil, nil
}
func (m *manualMockRepo) FetchAll(ctx context.Context) ([]models.IncidentLog, error) { return nil, nil }
func (m *manualMockRepo) Create(ctx context.Context, log *models.IncidentLog) error { return nil }
func (m *manualMockRepo) GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int) (*models.IncidentLog, error) {
	return nil, nil
}
func (m *manualMockRepo) UpdateStatus(ctx context.Context, tx *sql.Tx, id int, status string) error {
	return nil
}
func (m *manualMockRepo) SaveAISummary(ctx context.Context, summary *models.AISummary) error { return nil }
func (m *manualMockRepo) GetLatestSummary(ctx context.Context) (*models.AISummary, error)   { return nil, nil }

// ----------------------------------------------------------------------------
// Test dengan sqlmock untuk skenario database sebenarnya
// ----------------------------------------------------------------------------

func TestChangeLogStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, category_id, title, description, status, created_at, updated_at FROM incident_logs WHERE id = ? FOR UPDATE")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "category_id", "title", "description", "status", "created_at", "updated_at"}).
			AddRow(1, 1, 1, "Test", "Desc", "Menunggu", now, now))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE incident_logs SET status = ? WHERE id = ?")).
		WithArgs("Diproses", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := repository.NewMysqlLogRepository(db)
	uc := NewLogUsecase(repo, nil) // Perbaikan: tambahkan nil untuk eventBus

	err = uc.ChangeLogStatus(context.Background(), 1, "Diproses")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeLogStatus_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, category_id, title, description, status, created_at, updated_at FROM incident_logs WHERE id = ? FOR UPDATE")).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	repo := repository.NewMysqlLogRepository(db)
	uc := NewLogUsecase(repo, nil) // Perbaikan

	err = uc.ChangeLogStatus(context.Background(), 1, "Diproses")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeLogStatus_DeadlockRetry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	now := time.Now()
	// Percobaan 1: deadlock
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, category_id, title, description, status, created_at, updated_at FROM incident_logs WHERE id = ? FOR UPDATE")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "category_id", "title", "description", "status", "created_at", "updated_at"}).
			AddRow(1, 1, 1, "Test", "Desc", "Menunggu", now, now))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE incident_logs SET status = ? WHERE id = ?")).
		WithArgs("Selesai", 1).
		WillReturnError(errors.New("Deadlock found when trying to get lock; try restarting transaction"))
	mock.ExpectRollback()

	// Percobaan 2: sukses
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, category_id, title, description, status, created_at, updated_at FROM incident_logs WHERE id = ? FOR UPDATE")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "category_id", "title", "description", "status", "created_at", "updated_at"}).
			AddRow(1, 1, 1, "Test", "Desc", "Menunggu", now, now))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE incident_logs SET status = ? WHERE id = ?")).
		WithArgs("Selesai", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := repository.NewMysqlLogRepository(db)
	uc := NewLogUsecase(repo, nil) // Perbaikan

	err = uc.ChangeLogStatus(context.Background(), 1, "Selesai")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeLogStatus_DeadlockMaxRetriesExceeded(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	now := time.Now()
	for i := 0; i < 3; i++ {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, category_id, title, description, status, created_at, updated_at FROM incident_logs WHERE id = ? FOR UPDATE")).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "category_id", "title", "description", "status", "created_at", "updated_at"}).
				AddRow(1, 1, 1, "Test", "Desc", "Menunggu", now, now))
		mock.ExpectExec(regexp.QuoteMeta("UPDATE incident_logs SET status = ? WHERE id = ?")).
			WithArgs("Selesai", 1).
			WillReturnError(errors.New("Deadlock found when trying to get lock; try restarting transaction"))
		mock.ExpectRollback()
	}

	repo := repository.NewMysqlLogRepository(db)
	uc := NewLogUsecase(repo, nil) // Perbaikan

	err = uc.ChangeLogStatus(context.Background(), 1, "Selesai")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gagal setelah 3 percobaan")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ----------------------------------------------------------------------------
// Test context timeout menggunakan mock manual (tanpa sqlmock)
// ----------------------------------------------------------------------------
func TestChangeLogStatus_ContextTimeout(t *testing.T) {
	mockRepo := &manualMockRepo{
		beginTxFunc: func(ctx context.Context) (*sql.Tx, error) {
			return nil, context.Canceled
		},
	}
	uc := NewLogUsecase(mockRepo, nil) // Perbaikan
	err := uc.ChangeLogStatus(context.Background(), 1, "Diproses")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
}