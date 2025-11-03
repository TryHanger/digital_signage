package service

import (
	"fmt"
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/repository"
)

type TemplateService struct {
	repo *repository.TemplateRepository
}

func NewTemplateService(repo *repository.TemplateRepository) *TemplateService {
	return &TemplateService{repo: repo}
}

func (s *TemplateService) CreateTemplate(template *model.Template) error {
	if len(template.Blocks) == 0 {
		return fmt.Errorf("no template blocks")
	}

	for _, block := range template.Blocks {
		if block.StartTime.IsZero() || block.EndTime.IsZero() {
			return fmt.Errorf("start time or end time is empty")
		}
	}
	return s.repo.CreateTemplate(template)
}

func (s *TemplateService) GetAll() ([]model.Template, error) {
	return s.repo.GetAll()
}

func (s *TemplateService) GetByID(id uint) (*model.Template, error) {
	return s.repo.GetByID(id)
}

func (s *TemplateService) Delete(id uint) error {
	return s.repo.Delete(id)
}
