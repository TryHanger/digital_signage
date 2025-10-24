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

func (r *ScheduleRepository) DB() *gorm.DB {
	return r.db
}

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
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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

func (r *ScheduleRepository) GetSchedulesForDate(date time.Time) ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Preload("Content").
		Preload("Monitor").
		Preload("Location").
		Preload("Days").
		Joins("JOIN schedule_days ON schedules.id = schedule_days.schedule_id").
		Where("schedule_days.date = ?", date).
		Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) DeleteByID(id uint) (*model.Schedule, error) {
	var schedule model.Schedule

	tx := r.db.Begin()

	// Сначала находим, чтобы вернуть данные для удаления из кэша
	if err := tx.Preload("Days").First(&schedule, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Удаляем из таблицы schedule_days (если связь есть)
	if err := tx.Where("schedule_id = ?", id).Delete(&model.ScheduleDay{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Удаляем саму запись расписания
	if err := tx.Delete(&model.Schedule{}, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &schedule, tx.Commit().Error
}
