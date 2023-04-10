package model

import (
	"TrackMaster/pkg"
	"gorm.io/gorm"
	"time"
)

type Schedule struct {
	ProjectID string `gorm:"primary_key" json:"projectID"`

	// 定时任务的设置参数
	Interval time.Duration `gorm:"not null" json:"interval"`

	Name   string `json:"name"`
	Status bool   `gorm:"default:false" json:"status"`
}

func (s *Schedule) Create(db *gorm.DB) *pkg.Error {
	result := db.Create(s)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}

func (s *Schedule) Get(db *gorm.DB) *pkg.Error {
	result := db.First(&s)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}
