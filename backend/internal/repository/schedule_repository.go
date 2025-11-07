package repository

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"gorm.io/gorm"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(schedule *model.Schedule) error {
	return r.db.Create(schedule).Error
}

// GetAll возвращает все расписания с полным предзагрузом всех связанных данных
func (r *ScheduleRepository) GetAll() ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.
		Preload("Template.Blocks.Contents"). // контент шаблона
		Preload("Blocks.Contents").          // контент блоков расписания
		Preload("Monitor").
		Preload("MonitorGroup").
		Preload("Exceptions").
		Find(&schedules).Error
	return schedules, err
}

// GetByID возвращает конкретное расписание с полным preload
func (r *ScheduleRepository) GetByID(id uint) (*model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.
		Preload("Template.Blocks.Contents").
		Preload("Blocks.Contents").
		Preload("Monitor").
		Preload("MonitorGroup").
		Preload("Exceptions").
		First(&schedule, id).Error
	return &schedule, err
}

// Delete удаляет расписание и каскадно все блоки и контент
func (r *ScheduleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Schedule{}, id).Error
}

func (r *ScheduleRepository) WithTransaction(fn func(repo *ScheduleRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &ScheduleRepository{db: tx}
		return fn(txRepo)
	})
}

func (r *ScheduleRepository) GetTemplateWithBlocks(templateID uint) (*model.Template, error) {
	var template model.Template
	err := r.db.
		Preload("Blocks.Contents").
		First(&template, templateID).Error
	return &template, err
}

func (r *ScheduleRepository) CreateScheduleBlock(block *model.ScheduleBlock) error {
	return r.db.Create(block).Error
}

func (r *ScheduleRepository) CreateScheduleContent(content *model.ScheduleContent) error {
	return r.db.Create(content).Error
}
