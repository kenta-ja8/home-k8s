package repository

import (
	"time"

	"github.com/kenta-ja8/home-k8s-app/pkg/helper"
)

type BaseModel struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBaseModel() BaseModel {
	return BaseModel{
		ID:        helper.NewUUID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
