package model

import "time"

type Schedule struct {
	ID         uint          `json:"id" gorm:"primaryKey"`
	ContentID  uint          `json:"content_id" gorm:"not null"`
	MonitorID  *uint         `json:"monitor_id"`
	LocationID *uint         `json:"location_id"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Priority   int           `json:"priority" gorm:"default:0"`
	CreatedAt  time.Time     `json:"created_at"`
	Content    Content       `json:"content" gorm:"foreignKey:ContentID"`
	Monitor    *Monitor      `json:"monitor" gorm:"foreignKey:MonitorID"`
	Location   *Location     `json:"location" gorm:"foreignKey:LocationID"`
	Days       []ScheduleDay `json:"days"`
}

type ScheduleDay struct {
	ID         uint      `gorm:"primaryKey"`
	ScheduleID uint      `gorm:"not null;index"`
	Date       time.Time `gorm:"not null"`
}
