package models

import "time"

type AISummary struct {
	ID             int       `json:"id"`
	SummaryText    string    `json:"summary_text"`
	LogIDsAnalyzed string    `json:"log_ids_analyzed"` // comma separated, e.g. "1,2,3"
	CreatedAt      time.Time `json:"created_at"`
}