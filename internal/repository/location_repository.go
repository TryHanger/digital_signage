package repository

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"gorm.io/gorm"
)

type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Create(location *model.Location) error {
	return r.db.Create(location).Error
}

func (r *LocationRepository) GetAll() ([]model.Location, error) {
	var locations []model.Location
	err := r.db.Preload("Monitors").Find(&locations).Error
	return locations, err
}

func (r *LocationRepository) GetByID(id uint) (*model.Location, error) {
	var location model.Location
	err := r.db.Preload("Monitors").First(&location, id).Error
	return &location, err
}

func (r *LocationRepository) Update(location *model.Location) error {
	return r.db.Save(location).Error
}

func (r *LocationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Location{}, id).Error
}
