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
	Filter string `gorm:"unique" json:"filter"`

	Events    string     `json:"events"` // 存event的id数组
	EventLogs []EventLog `gorm:"foreignKey:RecordID" json:"eventLogs"`
	ProjectID string     `gorm:"index;not null" json:"projectID" binding:"required,max=32"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (r *Record) Create(db *gorm.DB) error {
	result := db.Create(r)
	return result.Error
}

func (r *Record) Get(db *gorm.DB) error {
	result := db.Preload("EventLogs.FieldLogs").First(&r)
	return result.Error
}

func (r *Record) Update(db *gorm.DB) error {
	err := r.Get(db)
	if err != nil {
		return err
	}

	result := db.Model(&r).Update("status", OFF)
	return result.Error
}
