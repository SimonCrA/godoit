package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Description      string
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
