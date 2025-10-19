package model

type Location struct {
	ID       uint `gorm:"primary_key"`
	Name     string
	Monitors []Monitor `gorm:"foreignKey:LocationID"`
}
