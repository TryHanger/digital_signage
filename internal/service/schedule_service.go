package service

import (
	"fmt"
	"github.com/TryHanger/digital_signage/internal/cache"
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/repository"
	"github.com/TryHanger/digital_signage/internal/socket"
	"gorm.io/gorm"
	"time"
)

type ScheduleService struct {
	repo     *repository.ScheduleRepository
	cache    *cache.ScheduleCache
	notifier *socket.WebSocketNotifier
}

func NewScheduleService(repo *repository.ScheduleRepository, cache *cache.ScheduleCache, notifier *socket.WebSocketNotifier) *ScheduleService {
	return &ScheduleService{repo: repo, cache: cache, notifier: notifier}
}

func (s *ScheduleService) CreateSchedule(schedule *model.Schedule) ([]model.Schedule, error) {
	// Проверяем конфликты перед созданием
	conflicts, err := s.repo.FindConflicts(schedule)
	if err != nil {
		return nil, err
	}

	if len(conflicts) > 0 {
		return conflicts, fmt.Errorf("conflict_with_existing")
	}

	// Транзакционно создаём новое расписание
	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.repo.Create(schedule); err != nil {
			return err
		}

		// Добавляем в кэш точечно
		s.cache.Add(*schedule)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Если расписание на сегодня — уведомляем мониторы
	if s.isToday(schedule) {
		//go s.notifier.BroadcastScheduleUpdate("created", schedule)
	}

	return nil, nil
}

func (s *ScheduleService) isToday(schedule *model.Schedule) bool {
	today := time.Now().Truncate(24 * time.Hour)
	for _, d := range schedule.Days {
		if d.Date.Truncate(24 * time.Hour).Equal(today) {
			return true
		}
	}
	return false
}

func (s *ScheduleService) GetAll() ([]model.Schedule, error) {
	return s.repo.GetAll()
}

func (s *ScheduleService) GetByID(id uint) (*model.Schedule, error) {
	return s.repo.GetByID(id)
}

func (s *ScheduleService) UpdateSchedules(schedules []model.Schedule) error {
	// Проверяем конфликты перед обновлением
	for _, sched := range schedules {
		conflicts, err := s.repo.FindConflicts(&sched)
		if err != nil {
			return err
		}
		if len(conflicts) > 0 {
			return fmt.Errorf("update_conflict_with_existing")
		}
	}

	// Обновляем в БД
	if err := s.repo.UpdateSchedules(schedules); err != nil {
		return err
	}

	// Обновляем только нужные элементы в кэше
	for _, sched := range schedules {
		s.cache.Update(sched)

		if s.isToday(&sched) {
			//go s.notifier.BroadcastScheduleUpdate("updated", &sched)
		}
	}

	return nil
}

func (s *ScheduleService) LoadDailyCache() error {
	today := time.Now().Truncate(24 * time.Hour)
	schedules, err := s.repo.GetSchedulesForDate(today)
	if err != nil {
		return err
	}
	s.cache.Set(schedules)
	return nil
}

func (s *ScheduleService) GetCachedSchedules() []model.Schedule {
	return s.cache.Get()
}

func (s *ScheduleService) SendSchedulesToMonitor(monitorID uint) {
	schedules := s.cache.GetByMonitorID(monitorID)
	for _, schedule := range schedules {
		s.notifier.NotifyScheduleUpdate(monitorID, schedule)
	}
}

func (s *ScheduleService) DeleteSchedule(id uint) error {
	// Удаляем из БД (можно через транзакцию)
	schedule, err := s.repo.DeleteByID(id)
	if err != nil {
		return err
	}

	// Удаляем из кэша
	s.cache.Delete(id)

	// Если удалённое расписание было на сегодня — уведомляем мониторы
	if s.isToday(schedule) {
		//go s.notifier.BroadcastScheduleUpdate("deleted", schedule)
		s.notifier.NotifyScheduleUpdate(*schedule.MonitorID, *schedule)
	}

	return nil
}
