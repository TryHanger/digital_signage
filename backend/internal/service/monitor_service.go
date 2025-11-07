package service

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/repository"
	"github.com/TryHanger/digital_signage/internal/utils"
	"strings"
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
	for {
		monitor.Token = utils.GenerateShortToken()
		err := s.repo.Create(monitor)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				// Повтор токена — сгенерируем заново
				continue
			}
			return err
		}
		break
	}
	return nil
}
