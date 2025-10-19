package service

import (
	"fmt"
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/repository"
)

type ScheduleService struct {
	repo *repository.ScheduleRepository
}

func NewScheduleService(repo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) CreateSchedule(schedule *model.Schedule) (*[]model.Schedule, error) {
	conflicts, err := s.repo.FindConflicts(schedule)
	if err != nil {
		return nil, err
	}

	if len(conflicts) > 0 {
		return &conflicts, fmt.Errorf("conflict_with_existing")
	}

	if err := s.repo.Create(schedule); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ScheduleService) GetAll() ([]model.Schedule, error) {
	return s.repo.GetAll()
}

func (s *ScheduleService) GetByID(id uint) (*model.Schedule, error) {
	return s.repo.GetByID(id)
}

func (s *ScheduleService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ScheduleService) UpdateSchedules(schedules []model.Schedule) error {
	// Здесь можно сделать проверки перед обновлением, например:
	// - не выходят ли новые времена за рамки дня
	// - не пересекаются ли с другими активными расписаниями
	// - корректность приоритетов (например, уникальность приоритета в один момент)

	return s.repo.UpdateSchedules(schedules)
}
