package cache

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/utils"
	"sync"
	"time"
)

type ScheduleCache struct {
	mu        sync.RWMutex
	Schedules []model.Schedule
}

func NewScheduleCache() *ScheduleCache {
	return &ScheduleCache{}
}

func (c *ScheduleCache) Set(schedules []model.Schedule) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Schedules = schedules
}

func (c *ScheduleCache) Get() []model.Schedule {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Schedules
}

func (c *ScheduleCache) Add(schedule model.Schedule) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Schedules = append(c.Schedules, schedule)
}

func (c *ScheduleCache) Update(updated model.Schedule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, sched := range c.Schedules {
		if sched.ID == updated.ID {
			c.Schedules[i] = updated
			return
		}
	}
}

func (c *ScheduleCache) Delete(id uint) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, sched := range c.Schedules {
		if sched.ID == id {
			// удаляем элемент с индексом i
			c.Schedules = append(c.Schedules[:i], c.Schedules[i+1:]...)
			return
		}
	}
}

func (c *ScheduleCache) GetAllToday() []model.Schedule {
	c.mu.RLock()
	defer c.mu.RUnlock()

	today := time.Now().Truncate(24 * time.Hour)
	var result []model.Schedule

	for _, sched := range c.Schedules {
		if utils.IsActiveToday(sched, today) {
			result = append(result, sched)
		}
	}

	return result
}

func (c *ScheduleCache) GetByMonitorID(monitorID uint) []model.Schedule {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []model.Schedule
	for _, s := range c.Schedules {
		if s.MonitorID != nil && *s.MonitorID == monitorID {
			result = append(result, s)
		}
	}
	return result
}
