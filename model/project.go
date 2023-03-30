package model

import (
	"TrackMaster/pkg"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	ID        string    `gorm:"primaryKey"`
	Name      string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewProject(id, name string) *Project {
	return &Project{
		ID:   id,
		Name: name,
	}
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

type SwaggerProject struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required,min=2,max=500"`
	CreatedAt time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updateAt"`
}

type SwaggerProjects struct {
	Data  []*SwaggerProject
	Pager *pkg.Pager
}
