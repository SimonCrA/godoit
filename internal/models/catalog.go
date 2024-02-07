package models

import "gorm.io/gorm"

type CatStatus struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Description   string `gorm:"not null"`
	LogicalDelete bool   `gorm:"default:false"`
}

type CatCategory struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Description   string `gorm:"not null"`
	LogicalDelete bool   `gorm:"default:false"`
}
