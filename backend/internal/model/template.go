package model

import "time"

type Template struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Blocks      []TemplateBlock `json:"blocks" gorm:"foreignKey:TemplateID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type TemplateBlock struct {
	ID         uint              `json:"id" gorm:"primaryKey"`
	TemplateID uint              `json:"templateId"`
	Name       string            `json:"name"`
	StartTime  time.Time         `json:"startTime"`
	EndTime    time.Time         `json:"endTime"`
	Contents   []TemplateContent `json:"contents" gorm:"foreignKey:BlockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TemplateContent struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	BlockID   uint   `json:"blockId"`
	ContentID uint   `json:"contentId"`
	Type      string `json:"type"`
	Order     int    `json:"order"`
	Duration  int    `json:"duration"`
}
