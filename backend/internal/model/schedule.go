package model

import (
	"github.com/lib/pq"
	"time"
)

// ========== ПОВТОРЕНИЯ ==========
type RepeatType string

const (
	RepeatNone   RepeatType = "none"
	RepeatDaily  RepeatType = "daily"
	RepeatWeekly RepeatType = "weekly"
)

// ========== ДНИ НЕДЕЛИ ==========
const (
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
	Thursday  = 4
	Friday    = 5
	Saturday  = 6
	Sunday    = 7
)

// ========== РАСПИСАНИЕ ==========
type Schedule struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description,omitempty"`

	// Шаблон (для копирования блоков)
	TemplateID uint      `json:"templateId" gorm:"not null"`
	Template   *Template `json:"template,omitempty"`

	// Устройства: либо Location, либо Group, либо отдельные Monitors
	LocationID *uint         `json:"locationId,omitempty"`
	Location   *Location     `json:"location,omitempty"`
	GroupID    *uint         `json:"groupId,omitempty"`
	Group      *MonitorGroup `json:"group,omitempty"`
	Monitors   []Monitor     `json:"monitors,omitempty" gorm:"many2many:schedule_monitors"`

	// Период
	StartDate time.Time  `json:"startDate" gorm:"not null"`
	EndDate   *time.Time `json:"endDate,omitempty"`

	// Повторение
	RepeatType RepeatType          `json:"repeatType" gorm:"default:'daily'"`
	Weekdays   pq.Int64Array       `json:"weekdays,omitempty" gorm:"type:integer[]"` // для weekly-repeat
	Exceptions []ScheduleException `json:"exceptions" gorm:"foreignKey:ScheduleID"`

	// Блоки
	Blocks []ScheduleBlock `json:"blocks" gorm:"foreignKey:ScheduleID"`

	// Статус
	IsActive bool `json:"isActive" gorm:"default:true"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ========== БЛОКИ РАСПИСАНИЯ ==========
type ScheduleBlock struct {
	ID         uint `json:"id" gorm:"primaryKey"`
	ScheduleID uint `json:"scheduleId" gorm:"not null"`

	// Скопировано из TemplateBlock
	Name      string `json:"name"`
	StartTime string `json:"startTime"` // "08:00"
	EndTime   string `json:"endTime"`   // "22:00"
	Position  int    `json:"position"`

	// Контент
	Items []ScheduleBlockItem `json:"items" gorm:"foreignKey:BlockID"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ScheduleBlockItem struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	BlockID   uint     `json:"blockId" gorm:"not null"`
	ContentID uint     `json:"contentId" gorm:"not null"`
	Content   *Content `json:"content,omitempty"`

	Position int  `json:"position"`
	Duration *int `json:"duration,omitempty"` // для фото

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ========== ИСКЛЮЧЕНИЯ ==========
type ScheduleException struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ScheduleID uint      `json:"scheduleId" gorm:"not null"`
	Date       time.Time `json:"date" gorm:"type:date;not null"`
	Reason     string    `json:"reason,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
