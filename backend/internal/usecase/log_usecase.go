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
	repo repository.LogRepository
}

func NewLogUsecase(repo repository.LogRepository) LogUsecase {
	return &logUsecase{repo: repo}
}

func (u *logUsecase) GetAllLogs(ctx context.Context) ([]models.IncidentLog, error) {
	return u.repo.FetchAll(ctx)
}

func (u *logUsecase) ReportIncident(ctx context.Context, log *models.IncidentLog) error {
	log.Status = "Menunggu"
	return u.repo.Create(ctx, log)
}

// Alur Transaksi Menggunakan Pessimistic Locking - DIPERBAIKI
func (u *logUsecase) ChangeLogStatus(ctx context.Context, id int, newStatus string) error {
	db := u.repo.GetDB() 
	if db == nil {
		return fmt.Errorf("koneksi database tidak tersedia")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Kunci baris data dengan FOR UPDATE
	_, err = u.repo.GetByIDForUpdate(ctx, tx, id)
	if err != nil {
		return err
	}

	// 2. Eksekusi update status data di dalam kunci pengaman
	err = u.repo.UpdateStatus(ctx, tx, id, newStatus)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Logika Perangkuman Massal Otomatis (Lazy Trigger) via Short-Polling API
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

	// Evaluasi: Jika belum ada summary, atau ada penambahan laporan baru yang belum dirangkum
	if currentSummary == nil || currentSummary.LogIDsAnalyzed != joinedIDs {
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return &models.AISummary{SummaryText: "Sistem AI Summary belum siap (API Key Kosong)."}, nil
		}

		prompt := fmt.Sprintf("Berikan ringkasan eksekutif singkat dalam Bahasa Indonesia mengenai kumpulan laporan gangguan operasional berikut, temukan pola kerusakan utamanya:\n%s", strings.Join(logTexts, "\n"))
		
		aiText, err := callGeminiAPI(prompt, apiKey)
		if err != nil {
			return currentSummary, nil // Kembalikan data lama sebagai fallback jika API gagal/limit
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

// Panggilan HTTP murni ke Gemini API tanpa library eksternal yang memberatkan performa
func callGeminiAPI(prompt, key string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", key)
	
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}
	
	bytesPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytesPayload))
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
		return "", fmt.Errorf("status kode tidak normal: %d", resp.StatusCode)
	}

	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return "", err
	}

	// Parsing JSON manual yang aman dari response Gemini API
	candidates, ok := responseMap["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("gagal membaca format response")
	}
	candidate := candidates[0].(map[string]interface{})
	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("format response tidak valid")
	}
	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("tidak ada parts dalam response")
	}
	part, ok := parts[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("format part tidak valid")
	}
	textResult, ok := part["text"].(string)
	if !ok {
		return "", fmt.Errorf("text tidak ditemukan dalam response")
	}

	return textResult, nil
}