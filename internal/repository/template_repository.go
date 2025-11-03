package repository

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"gorm.io/gorm"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) CreateTemplate(template *model.Template) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(template).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *TemplateRepository) GetAll() ([]model.Template, error) {
	var templates []model.Template
	err := r.db.Preload("Blocks").Preload("Blocks.Contents").Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) GetByID(id uint) (*model.Template, error) {
	var template model.Template
	err := r.db.Preload("Blocks").Preload("Blocks.Contents").First(&template, id).Error
	return &template, err
}

func (r *TemplateRepository) Delete(id uint) error {
	return r.db.Delete(&model.Template{}, id).Error
}
