package model

import "time"

type Content struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}
