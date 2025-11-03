package service

import (
	"time"

	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/utils"
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
