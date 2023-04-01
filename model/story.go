package model

import (
	"TrackMaster/pkg"
	"gorm.io/gorm"
	"time"
)

type Story struct {
	ProjectID   string    `gorm:"index;not null" json:"projectID" binding:"required,max=32"`
	ID          string    `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Name        string    `gorm:"not null" json:"name" binding:"required,min=2,max=50"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	Events      []Event   `gorm:"foreignKey:StoryID"`
}

func (s *Story) List(db *gorm.DB, pageOffset, pageSize int) ([]Story, int64, error) {
	var stories []Story
	result := db.Where("project_id = ?", s.ProjectID).Find(&stories)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if pageOffset >= 0 && pageSize > 0 {
		result = result.Offset(pageOffset).Limit(pageSize).Find(&stories)
	}

	// 将来要按条件过滤可以写在这里

	totalRow := result.RowsAffected
	return stories, totalRow, nil
}

func (s *Story) Get(db *gorm.DB) error {
	result := db.Preload("Events.Fields").Where("id = ?", s.ID).First(&s)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type Stories struct {
	Data  []*Story
	Pager *pkg.Pager
}
