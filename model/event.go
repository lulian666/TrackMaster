package model

import "time"

type Event struct {
	StoryID     string    `gorm:"index;not null" json:"storyID" binding:"required,max=32"`
	ID          string    `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Name        string    `gorm:"not null" json:"name" binding:"required,min=2,max=50"`
	Description string    `json:"description"`
	OnTrail     bool      `gorm:"default:false" json:"onTrail"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	Fields      []Field   `gorm:"foreignKey:EventID"`
}
