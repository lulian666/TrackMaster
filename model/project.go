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

func (p *Project) Create(db *gorm.DB) *pkg.Error {
	result := db.Create(p)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}

func (p *Project) List(db *gorm.DB) ([]Project, *pkg.Error) {
	var projects []Project
	result := db.Find(&projects)
	if result.Error != nil {
		return nil, pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return projects, nil
}

func (p *Project) Get(db *gorm.DB) *pkg.Error {
	result := db.First(&p)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}

type Projects struct {
	Data  []*Project
	Pager *pkg.Pager
}
