package repository

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"gorm.io/gorm"
	"time"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(s *model.Schedule) error {
	return r.db.Create(s).Error
}

func (r *ScheduleRepository) GetAll() ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Preload("Monitor").Preload("Content").Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) GetByID(id uint) (*model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.Preload("Monitor").Preload("Content").First(&schedule, id).Error
	return &schedule, err
}

func (r *ScheduleRepository) Update(s *model.Schedule) error {
	return r.db.Save(s).Error
}

func (r *ScheduleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Schedule{}, id).Error
}

func (r *ScheduleRepository) GetActiveByMonitorID(monitorID uint, now time.Time) (*model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.Preload("Content").
		Where("monitor_id = ? AND start_time <= ? AND end_time >= ?", monitorID, now, now).
		First(&schedule).Error
	return &schedule, err
}

func (r *ScheduleRepository) GetAllActive() ([]model.Schedule, error) {
	now := time.Now()
	var schedules []model.Schedule
	err := r.db.Preload("Content").Preload("Monitor").
		Where("start_time <= ? AND end_time >= ?", now, now).
		Find(&schedules).Error
	return schedules, err
}
