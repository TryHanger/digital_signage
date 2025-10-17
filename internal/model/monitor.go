package model

import "time"

type Monitor struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Location  string
	Status    string
	CreatedAt time.Time
}
