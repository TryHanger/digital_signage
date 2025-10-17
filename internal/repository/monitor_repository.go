package repository

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"gorm.io/gorm"
)

type MonitorRepository struct {
	db *gorm.DB
}

func NewMonitorRepository(db *gorm.DB) *MonitorRepository {
	return &MonitorRepository{db: db}
}

func (r *MonitorRepository) GetAll() ([]model.Monitor, error) {
	var monitors []model.Monitor
	err := r.db.Find(&monitors).Error
	if err != nil {
		return nil, err
	}
	return monitors, nil
}

func (r *MonitorRepository) Create(monitor *model.Monitor) error {
	return r.db.Create(monitor).Error
}
