package model

import (
	"TrackMaster/pkg"
	"database/sql/driver"
	"gorm.io/gorm"
	"time"
)

type status string

const (
	ON  status = "ON"
	OFF status = "OFF"
)

func (s *status) Scan(value interface{}) error {
	*s = status(value.([]byte))
	return nil
}

func (s *status) Value() (driver.Value, error) {
	return string(*s), nil
}

type Record struct {
	ID     string `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Name   string `gorm:"not null" json:"name" binding:"required,min=2,max=50"`
	Status status `sql:"type:ENUM('ON', 'OFF')"  json:"status"`
	Filter string `gorm:"unique" json:"filter"`

	//Events    string     `json:"events"` // 存event的id数组
	Events []Event `gorm:"many2many:record_events;"`

	EventLogs []EventLog `gorm:"foreignKey:RecordID" json:"eventLogs"`
	ProjectID string     `gorm:"index;not null" json:"projectID" binding:"required,max=32"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (r *Record) Create(db *gorm.DB) *pkg.Error {
	result := db.Create(r)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}

func (r *Record) Get(db *gorm.DB) *pkg.Error {
	result := db.Preload("EventLogs.FieldLogs").Preload("Events").First(&r)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}

func (r *Record) Update(db *gorm.DB) *pkg.Error {
	err := r.Get(db)
	if err != nil {
		return err
	}

	result := db.Model(&r).Update("status", OFF)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}
