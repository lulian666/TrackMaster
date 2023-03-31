package model

import (
	"gorm.io/gorm"
	"time"
)

type EnumValue struct {
	TypeId    string    `gorm:"index;not null" json:"typeID" binding:"required,max=32"`
	ID        string    `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Name      string    `gorm:"not null" json:"name" binding:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Type struct {
	ProjectID string    `gorm:"index;not null" json:"projectID" binding:"required,max=32"`
	ID        string    `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Type      string    `gorm:"not null" json:"type" binding:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (t *Type) Create(db *gorm.DB) error {
	result := db.Create(t)
	return result.Error
}

func (e *EnumValue) Create(db *gorm.DB) error {
	result := db.Create(e)
	return result.Error
}
