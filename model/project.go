package model

import (
	"TrackMaster/pkg"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name" binding:"required,min=2,max=500"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (p *Project) Create(db *gorm.DB) error {
	result := db.Create(p)
	return result.Error
}

func (p *Project) List(db *gorm.DB) ([]Project, error) {
	var projects []Project
	var err error
	result := db.Find(&projects)
	if result.Error != nil {
		return nil, err
	}
	return projects, nil
}

type Projects struct {
	Data  []*Project
	Pager *pkg.Pager
}
