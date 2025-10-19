package model

import "time"

type Monitor struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	Status     string
	CreatedAt  time.Time
	LocationID uint
	Location   *Location `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
