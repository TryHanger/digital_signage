package model

import "time"

type Schedule struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	MonitorID uint      `gorm:"primary_key" json:"monitor_id"`
	ContentID uint      `gorm:"primary_key" json:"content_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	Monitor Monitor `gorm:"foreignKey:MonitorID" json:"monitor"`
	Content Content `gorm:"foreignKey:ContentID" json:"content""`
}
