package model

import "time"

type Monitor struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Token      string `gorm:"uniqueIndex;size:32;not null"`
	Status     string
	CreatedAt  time.Time
	LocationID uint
	Location   *Location `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
