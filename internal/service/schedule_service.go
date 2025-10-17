package service

import (
	"time"

	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/repository"
)

type ScheduleService struct {
	repo *repository.ScheduleRepository
}

func NewScheduleService(repo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) Create(schedule *model.Schedule) error {
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

func (s *ScheduleService) GetActiveContent(monitorID uint) (*model.Schedule, error) {
	now := time.Now()
	return s.repo.GetActiveByMonitorID(monitorID, now)
}

func (s *ScheduleService) GetAllActive() ([]model.Schedule, error) {
	return s.repo.GetAllActive()
}
