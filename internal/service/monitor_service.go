package service

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/repository"
)

type MonitorService struct {
	repo *repository.MonitorRepository
}

func NewMonitorService(repo *repository.MonitorRepository) *MonitorService {
	return &MonitorService{repo: repo}
}

func (s *MonitorService) GetAllMonitors() ([]model.Monitor, error) {
	return s.repo.GetAll()
}

func (s *MonitorService) CreateMonitor(monitor *model.Monitor) error {
	return s.repo.Create(monitor)
}
