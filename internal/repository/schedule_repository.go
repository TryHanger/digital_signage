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

// Создание расписания и связанных дней
//func (r *ScheduleRepository) Create(schedule *model.Schedule) error {
//	return r.db.Create(schedule).Error
//}

func (r *ScheduleRepository) Create(schedule *model.Schedule) error {
	if err := r.db.Create(schedule).Error; err != nil {
		return err
	}

	// После создания — загрузим связи
	return r.db.
		Preload("Content").
		Preload("Monitor").
		Preload("Location").
		Preload("Days").
		First(schedule, schedule.ID).Error
}

// Получение всех расписаний
func (r *ScheduleRepository) GetAll() ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Preload("Days").Preload("Content").Find(&schedules).Error
	return schedules, err
}

// Получить одно по ID
func (r *ScheduleRepository) GetByID(id uint) (*model.Schedule, error) {
	var schedule model.Schedule
	if err := r.db.Preload("Content").
		Preload("Monitor").
		Preload("Location").
		Preload("Days").
		First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

// Удаление
func (r *ScheduleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Schedule{}, id).Error
}

func (r *ScheduleRepository) FindConflicts(schedule *model.Schedule) ([]model.Schedule, error) {
	var conflicts []model.Schedule

	query := r.db.
		Preload("Content").
		Preload("Location").
		Preload("Days")

	if schedule.MonitorID != nil {
		query = query.Where("monitor_id = ?", schedule.MonitorID)
	} else if schedule.LocationID != nil {
		query = query.Where("location_id = ?", schedule.LocationID)
	}

	for _, d := range schedule.Days {
		var dayConflicts []model.Schedule
		query.
			Joins("JOIN schedule_days ON schedules.id = schedule_days.schedule_id").
			Where("schedule_days.date = ?", d.Date).
			Where("schedules.start_time < ? AND schedules.end_time > ?", schedule.EndTime, schedule.StartTime).
			Find(&dayConflicts)

		conflicts = append(conflicts, dayConflicts...)
	}

	return conflicts, nil
}

func (r *ScheduleRepository) UpdateSchedules(schedules []model.Schedule) error {
	tx := r.db.Begin()

	for _, sched := range schedules {
		if err := tx.Model(&model.Schedule{}).
			Where("id = ?", sched.ID).
			Updates(map[string]interface{}{
				"priority":   sched.Priority,
				"start_time": sched.StartTime,
				"end_time":   sched.EndTime,
			}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
