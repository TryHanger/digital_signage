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

func (s *ScheduleService) Create(schedule *model.Schedule) error {
	if schedule.Name == "" {
		return errors.New("name is required")
	}
	if schedule.TemplateID == 0 {
		return errors.New("template is required")
	}
	return s.repo.Create(schedule)
}

func (s *ScheduleService) GetAll() ([]model.Schedule, error) {
	return s.repo.GetAll()
}

func (s *ScheduleService) GetByID(id uint) (*model.Schedule, error) {
	return s.repo.GetByID(id)
}

func (s *ScheduleService) Update(schedule *model.Schedule) error {
	return s.repo.Update(schedule)
}

func (s *ScheduleService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ScheduleService) GetActiveOn(date time.Time) ([]model.Schedule, error) {
	return s.repo.GetActiveOn(date)
}
