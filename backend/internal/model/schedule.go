package model

import (
	"github.com/lib/pq"
	"time"
)

type Schedule struct {
	ID             uint                `json:"id" gorm:"primaryKey"`
	Name           string              `json:"name"`
	TemplateID     uint                `json:"templateId"`
	Template       *Template           `json:"template" gorm:"constraint:OnDelete:CASCADE"`
	MonitorID      *uint               `json:"monitorId"`
	Monitor        *Monitor            `json:"monitor" gorm:"constraint:OnDelete:CASCADE"`
	MonitorGroupID *uint               `json:"monitorGroupId"`
	MonitorGroup   *MonitorGroup       `json:"monitorGroup" gorm:"constraint:OnDelete:CASCADE"`
	DateStart      time.Time           `json:"dateStart"`
	DateEnd        time.Time           `json:"dateEnd"`
	RepeatPattern  string              `json:"repeatPattern" gorm:"default:none"` // "none", "daily", "weekly", "monthly"
	DaysOfWeek     pq.Int64Array       `json:"daysOfWeek" gorm:"type:integer[]"`  // если RepeatPattern = weekly
	Mode           string              `json:"mode" gorm:"default:rotation"`      // "rotation" или "override"
	Priority       int                 `json:"priority" gorm:"default:1"`
	Blocks         []ScheduleBlock     `json:"blocks" gorm:"foreignKey:ScheduleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Exceptions     []ScheduleException `json:"exceptions" gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`
}

type ScheduleBlock struct {
	ID         uint              `json:"id" gorm:"primaryKey"`
	ScheduleID uint              `json:"scheduleId"` // к какому расписанию относится
	BlockID    uint              `json:"blockId"`    // от какого шаблона скопирован
	Name       string            `json:"name"`
	StartTime  time.Time         `json:"startTime"`
	EndTime    time.Time         `json:"endTime"`
	Contents   []ScheduleContent `json:"contents" gorm:"foreignKey:ScheduleBlockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

type ScheduleContent struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	ScheduleBlockID uint   `json:"scheduleBlockId"` // принадлежит ScheduleBlock
	ContentID       uint   `json:"contentId"`       // оригинальный контент
	Type            string `json:"type"`
	Order           int    `json:"order"`    // порядок показа
	Duration        int    `json:"duration"` // длительность показа
}
type ScheduleException struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ScheduleID uint      `json:"scheduleId"`
	Date       time.Time `json:"date"`
}
