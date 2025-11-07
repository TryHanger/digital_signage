package service

import (
	"github.com/TryHanger/digital_signage/backend/internal/model"
	"github.com/TryHanger/digital_signage/backend/internal/repository"
)

type LocationService struct {
	repo *repository.LocationRepository
}

func NewLocationService(repo *repository.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) Create(location *model.Location) error {
	return s.repo.Create(location)
}

func (s *LocationService) GetAll() ([]model.Location, error) {
	return s.repo.GetAll()
}

func (s *LocationService) GetByID(id uint) (*model.Location, error) {
	return s.repo.GetByID(id)
}

func (s *LocationService) Update(location *model.Location) error {
	return s.repo.Update(location)
}

func (s *LocationService) Delete(id uint) error {
	return s.repo.Delete(id)
}
