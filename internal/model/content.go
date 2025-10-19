package model

import "time"

type Content struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Type        string
	Path        string
	Description string
	Duration    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
