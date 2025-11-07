package model

type Location struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Name     string    `json:"name"`
	Monitors []Monitor `json:"monitors" gorm:"foreignKey:LocationID"`
}
