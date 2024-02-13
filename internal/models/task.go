package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title            string    `gorm:"size:30" validate:"min=3,max=30"`
	Description      string    `gorm:"size:255" validate:"required,min=5,max=255"`
	CurrentTaskDate  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	LastStatusChange time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	LogicalDelete    bool      `gorm:"default:false"`
	FkIdUser         int
	User             User `gorm:"foreignKey:FkIdUser"`
	FkIdCatStatus    int
	CatStatus        CatStatus `gorm:"foreignKey:FkIdCatStatus"`
	FkIdCatCategory  int
	CatCategory      CatCategory `gorm:"foreignKey:FkIdCatCategory"`
}

type TaskResponse struct {
	IdTask uint `json:"idTask"`
}
