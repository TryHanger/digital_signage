package model

import "time"

type Schedule struct {
	ID         uint          `json:"id" gorm:"primaryKey"`
	ContentID  uint          `json:"contentID" gorm:"not null"`
	MonitorID  *uint         `json:"monitorID"`
	LocationID *uint         `json:"locationID"`
	StartTime  time.Time     `json:"startTime"`
	EndTime    time.Time     `json:"endTime"`
	Priority   int           `json:"priority" gorm:"default:0"`
	CreatedAt  time.Time     `json:"createdAt"`
	Content    Content       `json:"content" gorm:"foreignKey:ContentID"`
	Monitor    *Monitor      `json:"monitor" gorm:"foreignKey:MonitorID"`
	Location   *Location     `json:"location" gorm:"foreignKey:LocationID"`
	Days       []ScheduleDay `json:"days"`
}

type ScheduleDay struct {
	ID         uint      `gorm:"primaryKey"`
	ScheduleID uint      `gorm:"not null;index"`
	Date       time.Time `json:"date" gorm:"not null"`
}
