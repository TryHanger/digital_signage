package service

import (
	"errors"
	"github.com/TryHanger/digital_signage/backend/internal/model"
	"github.com/TryHanger/digital_signage/backend/internal/repository"
	"time"
)

type ScheduleService struct {
	repo *repository.ScheduleRepository
}

func NewScheduleService(repo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) CreateSchedule(schedule *model.Schedule) error {
	// ✅ Валидация
	if schedule.TemplateID == 0 {
		return errors.New("templateID is required")
	}
	if schedule.MonitorID == nil && schedule.MonitorGroupID == nil {
		return errors.New("schedule must be assigned to either monitor or group")
	}
	if schedule.DateStart.After(schedule.DateEnd) {
		return errors.New("dateStart cannot be after dateEnd")
	}
	if schedule.Mode != "rotation" && schedule.Mode != "override" {
		return errors.New("invalid mode")
	}

	if schedule.RepeatPattern == "" {
		schedule.RepeatPattern = "none"
	}

	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	// ✅ 1. Сохраняем расписание
	err := s.repo.Create(schedule)
	if err != nil {
		return err
	}

	// ✅ 2. Получаем шаблон с блоками
	template, err := s.repo.GetTemplateWithBlocks(schedule.TemplateID)
	if err != nil {
		return err
	}

	// ✅ 3. Копируем блоки шаблона в расписание
	for _, tmplBlock := range template.Blocks {
		scheduleBlock := model.ScheduleBlock{
			ScheduleID: schedule.ID,
			Name:       tmplBlock.Name,
			StartTime:  tmplBlock.StartTime,
			EndTime:    tmplBlock.EndTime,
		}

		if err := s.repo.CreateScheduleBlock(&scheduleBlock); err != nil {
			return err
		}

		// ✅ 4. Копируем контенты из шаблонного блока
		for _, tmplContent := range tmplBlock.Contents {
			scheduleContent := model.ScheduleContent{
				ScheduleBlockID: scheduleBlock.ID,
				ContentID:       tmplContent.ContentID,
				Type:            tmplContent.Type,
				Order:           tmplContent.Order,
				Duration:        tmplContent.Duration,
			}
			if err := s.repo.CreateScheduleContent(&scheduleContent); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ScheduleService) GetAllSchedules() ([]model.Schedule, error) {
	return s.repo.GetAll()
}

func (s *ScheduleService) GetScheduleByID(id uint) (*model.Schedule, error) {
	return s.repo.GetByID(id)
}

func (s *ScheduleService) DeleteSchedule(id uint) error {
	return s.repo.Delete(id)
}

// validateSchedule проверяет корректность основных полей
func validateSchedule(schedule *model.Schedule) error {
	if schedule.TemplateID == 0 {
		return errors.New("templateID is required")
	}
	if schedule.MonitorID == nil && schedule.MonitorGroupID == nil {
		return errors.New("schedule must be assigned to a monitor or group")
	}
	if schedule.DateStart.After(schedule.DateEnd) {
		return errors.New("dateStart cannot be after dateEnd")
	}
	if schedule.Mode != "" && schedule.Mode != "rotation" && schedule.Mode != "override" {
		return errors.New("invalid mode: must be 'rotation' or 'override'")
	}
	return nil
}
