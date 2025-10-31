package model

import "time"

type Content struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Type        string    `json:"type"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
