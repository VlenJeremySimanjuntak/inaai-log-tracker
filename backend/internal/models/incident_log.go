package models

import "time"

type IncidentLog struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	CategoryID  int       `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // Menunggu, Diproses, Selesai
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}