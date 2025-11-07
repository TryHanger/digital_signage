package service

import (
	"fmt"
	"sort"

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
	// validate non-overlapping blocks (compare time-of-day)
	if overlapsTemplateBlocks(template.Blocks) {
		return fmt.Errorf("template blocks have overlapping time ranges")
	}
	return s.repo.CreateTemplate(template)
}

func (s *TemplateService) UpdateTemplate(template *model.Template) error {
	if template.ID == 0 {
		return fmt.Errorf("template id is required")
	}
	if len(template.Blocks) == 0 {
		return fmt.Errorf("no template blocks")
	}

	for _, block := range template.Blocks {
		if block.StartTime.IsZero() || block.EndTime.IsZero() {
			return fmt.Errorf("start time or end time is empty")
		}
	}

	if overlapsTemplateBlocks(template.Blocks) {
		return fmt.Errorf("template blocks have overlapping time ranges")
	}

	return s.repo.UpdateTemplate(template)
}

// overlapsTemplateBlocks checks whether any TemplateBlock time ranges overlap (by time-of-day)
func overlapsTemplateBlocks(blocks []model.TemplateBlock) bool {
	type interval struct{ start, end int }
	var ivs []interval
	for _, b := range blocks {
		sh, sm := b.StartTime.Hour(), b.StartTime.Minute()
		eh, em := b.EndTime.Hour(), b.EndTime.Minute()
		start := sh*60 + sm
		end := eh*60 + em
		ivs = append(ivs, interval{start: start, end: end})
	}
	// sort by start
	sort.Slice(ivs, func(i, j int) bool { return ivs[i].start < ivs[j].start })
	for i := 0; i+1 < len(ivs); i++ {
		if ivs[i].end > ivs[i+1].start { // overlap
			return true
		}
	}
	return false
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
