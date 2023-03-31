package model

import "time"

type Field struct {
	EventID     string    `gorm:"index;not null" json:"eventID" binding:"required,max=32"`
	ID          string    `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Type        string    `json:"type"`
	TypeID      string    `json:"typeID"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
