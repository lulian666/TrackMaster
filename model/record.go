package model

import (
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
	Filter string `json:"filter"`

	EventLogs []EventLog `gorm:"foreignKey:RecordID" json:"eventLogs"`
	ProjectID string     `gorm:"index;not null" json:"projectID" binding:"required,max=32"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (r *Record) Create(db *gorm.DB) error {
	result := db.Create(r)
	return result.Error
}
