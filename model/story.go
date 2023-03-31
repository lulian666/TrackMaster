package model

import "time"

type Story struct {
	ProjectID   string    `gorm:"index;not null" json:"projectID" binding:"required,max=32"`
	ID          string    `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Name        string    `gorm:"not null" json:"name" binding:"required,min=2,max=50"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
