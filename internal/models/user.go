package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name                   string    `gorm:"size:255;not null" validate:"required,min=3"`
	Email                  string    `gorm:"uniqueIndex;not null" validate:"required,email"`
	Password               string    `gorm:"not null" validate:"required,min=6"`
	PasswordExpirationDate time.Time `gorm:"not null" `
	LastSession            time.Time `gorm:"not null" `
	LogicalDelete          bool      `gorm:"default:false"`
	FkIdCatStatus          int
	CatStatus              CatStatus `gorm:"foreignKey:FkIdCatStatus"`
}
