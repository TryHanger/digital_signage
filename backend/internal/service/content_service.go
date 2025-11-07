package service

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/repository"
)

type ContentService struct {
	repo *repository.ContentRepository
}

func NewContentService(repo *repository.ContentRepository) *ContentService {
	return &ContentService{repo: repo}
}

func (s *ContentService) Create(content *model.Content) error {
	return s.repo.Create(content)
}

func (s *ContentService) GetAll() ([]model.Content, error) {
	return s.repo.GetAll()
}

func (s *ContentService) GetByID(id uint) (*model.Content, error) {
	return s.repo.GetByID(id)
}

func (s *ContentService) Update(content *model.Content) error {
	return s.repo.Update(content)
}

func (s *ContentService) Delete(id uint) error {
	return s.repo.Delete(id)
}
