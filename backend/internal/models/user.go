package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"` // Teknisi / Admin
	CreatedAt time.Time `json:"created_at"`
}