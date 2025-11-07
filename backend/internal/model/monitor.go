package model

import "time"

type Monitor struct {
	ID         uint          `json:"id" gorm:"primaryKey"`
	Name       string        `json:"name" gorm:"not null"`
	Token      string        `json:"token" gorm:"uniqueIndex;size:32;not null"`
	Status     string        `json:"status"`
	CreatedAt  time.Time     `json:"createdAt"`
	LocationID uint          `json:"locationID"`
	Location   *Location     `json:"location" gorm:"constraint:OnDelete:CASCADE"`
	GroupID    *uint         `json:"groupID"`
	Group      *MonitorGroup `json:"group" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type MonitorGroup struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Name     string    `json:"name"`
	Monitors []Monitor `json:"monitors" gorm:"foreignKey:GroupID"`
}
