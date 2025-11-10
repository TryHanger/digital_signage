package repository

import (
	"github.com/TryHanger/digital_signage/backend/internal/model"
	"gorm.io/gorm"
	"time"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(schedule *model.Schedule) error {
	return r.db.Create(&schedule).Error
}

func (r *ScheduleRepository) GetAll() ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Preload("Blocks.Items.Content").
		Preload("Monitors").
		Preload("Location").
		Preload("Group").
		Preload("Exceptions").
		Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) GetByID(id uint) (*model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.Preload("Blocks.Items.Content").
		Preload("Monitors").
		Preload("Location").
		Preload("Group").
		Preload("Exceptions").
		First(&schedule, id).Error
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *ScheduleRepository) Update(schedule *model.Schedule) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(schedule).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *ScheduleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Schedule{}, id).Error
}

func (r *ScheduleRepository) GetActiveOn(date time.Time) ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Preload("Blocks.Items.Content").
		Where("start_date <= ?", date).
		Where("end_date IS NULL OR end_date >= ?", date).
		Find(&schedules).Error
	return schedules, err
}
