package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"backend-tracker/internal/event"
	"backend-tracker/internal/models"
	"backend-tracker/internal/repository"
)

type LogUsecase interface {
	GetAllLogs(ctx context.Context) ([]models.IncidentLog, error)
	ReportIncident(ctx context.Context, log *models.IncidentLog) error
	ChangeLogStatus(ctx context.Context, id int, newStatus string) error
	GetOrTriggerAggregateSummary(ctx context.Context) (*models.AISummary, error)
}

type logUsecase struct {
	repo     repository.LogRepository
	eventBus *event.Bus
}

func NewLogUsecase(repo repository.LogRepository, eventBus *event.Bus) LogUsecase {
	return &logUsecase{repo: repo, eventBus: eventBus}
}

func (u *logUsecase) GetAllLogs(ctx context.Context) ([]models.IncidentLog, error) {
	return u.repo.FetchAll(ctx)
}

func (u *logUsecase) ReportIncident(ctx context.Context, log *models.IncidentLog) error {
	log.Status = "Menunggu"
	err := u.repo.Create(ctx, log)
	if err == nil {
		// publish event for new incident
		u.eventBus.Publish(ctx, event.Event{
			Type: "log.created",
			Data: map[string]interface{}{
				"log_id": log.ID,
				"title":  log.Title,
			},
		})
	}
	return err
}

// changeStatusTx melakukan transaksi dengan pessimistic locking
func (u *logUsecase) changeStatusTx(ctx context.Context, id int, newStatus string) error {
	tx, err := u.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Kunci baris dengan FOR UPDATE
	_, err = u.repo.GetByIDForUpdate(ctx, tx, id)
	if err != nil {
		return err
	}

	// Update status
	err = u.repo.UpdateStatus(ctx, tx, id, newStatus)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ChangeLogStatus dengan retry deadlock dan publish event setelah sukses
func (u *logUsecase) ChangeLogStatus(ctx context.Context, id int, newStatus string) error {
	const maxRetries = 3
	baseDelay := 50 * time.Millisecond

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := u.changeStatusTx(ctx, id, newStatus)
		if err == nil {
			// Publikasikan event setelah sukses update
			if u.eventBus != nil {
				u.eventBus.Publish(ctx, event.Event{
					Type: "log.status.updated",
					Data: map[string]interface{}{
						"log_id": id,
						"status": newStatus,
					},
				})
			}
			return nil
		}

		if strings.Contains(err.Error(), "Deadlock") || strings.Contains(err.Error(), "1213") {
			delay := baseDelay * time.Duration(1<<attempt)
			time.Sleep(delay)
			lastErr = err
			continue
		}
		return err
	}
	return fmt.Errorf("gagal setelah %d percobaan: %w", maxRetries, lastErr)
}

// changeStatusTxWithEvent melakukan transaksi dan mengembalikan old status
func (u *logUsecase) changeStatusTxWithEvent(ctx context.Context, id int, newStatus string, oldStatusPtr *string) error {
	tx, err := u.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Kunci baris dengan FOR UPDATE dan ambil data lama
	log, err := u.repo.GetByIDForUpdate(ctx, tx, id)
	if err != nil {
		return err
	}
	*oldStatusPtr = log.Status

	// Update status
	err = u.repo.UpdateStatus(ctx, tx, id, newStatus)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetOrTriggerAggregateSummary (tidak berubah, tapi bisa subscribe event juga)
func (u *logUsecase) GetOrTriggerAggregateSummary(ctx context.Context) (*models.AISummary, error) {
	currentSummary, err := u.repo.GetLatestSummary(ctx)
	if err != nil {
		return nil, err
	}

	allLogs, err := u.repo.FetchAll(ctx)
	if err != nil || len(allLogs) == 0 {
		return currentSummary, err
	}

	var latestIDs []string
	var logTexts []string
	for _, l := range allLogs {
		latestIDs = append(latestIDs, strconv.Itoa(l.ID))
		logTexts = append(logTexts, fmt.Sprintf("[%s] %s: %s", l.Status, l.Title, l.Description))
	}
	joinedIDs := strings.Join(latestIDs, ",")

	if currentSummary == nil || currentSummary.LogIDsAnalyzed != joinedIDs {
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return &models.AISummary{SummaryText: "Sistem AI Summary belum siap (API Key Kosong)."}, nil
		}
		prompt := fmt.Sprintf("Berikan ringkasan eksekutif singkat dalam Bahasa Indonesia mengenai kumpulan laporan gangguan operasional berikut, temukan pola kerusakan utamanya:\n%s", strings.Join(logTexts, "\n"))
		aiText, err := callGeminiAPIWithContext(ctx, prompt, apiKey)
		if err != nil {
			return currentSummary, nil
		}
		newSummary := &models.AISummary{
			SummaryText:    aiText,
			LogIDsAnalyzed: joinedIDs,
		}
		_ = u.repo.SaveAISummary(ctx, newSummary)
		return newSummary, nil
	}
	return currentSummary, nil
}

func callGeminiAPIWithContext(ctx context.Context, prompt, key string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", key)
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	}
	bytesPayload, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bytesPayload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return "", err
	}
	candidates, ok := responseMap["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("no candidates")
	}
	candidate := candidates[0].(map[string]interface{})
	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid content")
	}
	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("no parts")
	}
	part, ok := parts[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid part")
	}
	text, ok := part["text"].(string)
	if !ok {
		return "", fmt.Errorf("no text")
	}
	return text, nil
}