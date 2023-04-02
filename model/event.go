package model

import (
	"gorm.io/gorm"
	"time"
)

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

func (e *Event) List(db *gorm.DB, eventLogIDs []string) ([]Event, int64, error) {
	var events []Event
	result := db.Model(Event{}).Preload("Fields").Where("id in (?)", eventLogIDs).Find(&events)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	totalRow := result.RowsAffected
	return events, totalRow, nil
}

func (e *Event) ListEventName(db *gorm.DB, eventIDs []string) ([]string, int64, error) {
	var events []string
	result := db.Model(e).Where("id in (?)", eventIDs).Pluck("name", &events)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	totalRow := result.RowsAffected
	return events, totalRow, nil
}
