package model

import "time"

type Monitor struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	Token      string    `json:"token" gorm:"uniqueIndex;size:32;not null"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	LocationID uint      `json:"locationID"`
	Location   *Location `json:"location" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
