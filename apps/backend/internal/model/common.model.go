package model

import (
	"gorm.io/gorm"
	"time"
)

type (
	CommonModel struct {
		ID        uint           `gorm:"primarykey" json:"id"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}
)
