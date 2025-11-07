package repository

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"gorm.io/gorm"
)

type ContentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

func (r *ContentRepository) Create(content *model.Content) error {
	return r.db.Create(content).Error
}

func (r *ContentRepository) GetAll() ([]model.Content, error) {
	var contents []model.Content
	err := r.db.Find(&contents).Error
	return contents, err
}

func (r *ContentRepository) GetByID(id uint) (*model.Content, error) {
	var content model.Content
	err := r.db.First(&content, id).Error
	return &content, err
}

func (r *ContentRepository) Update(content *model.Content) error {
	return r.db.Save(content).Error
}

func (r *ContentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Content{}, id).Error
}
