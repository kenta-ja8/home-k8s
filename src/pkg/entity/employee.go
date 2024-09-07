package entity

import "time"

type Employee struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
