package service

import (
	"github.com/TryHanger/digital_signage/backend/internal/model"
	"github.com/TryHanger/digital_signage/backend/internal/utils"
	"time"
)

func (s *ScheduleService) LoadDailyCache() error {
	schedules, err := s.repo.GetAll() // получаем все расписания из БД
	if err != nil {
		return err
	}

	activeToday := []model.Schedule{}
	today := time.Now().Truncate(24 * time.Hour)

	for _, sched := range schedules {
		if utils.IsActiveToday(sched, today) {
			activeToday = append(activeToday, sched)
		}
	}

	// Обновляем кэш
	// s.cache.Set(activeToday)
	return nil
}
