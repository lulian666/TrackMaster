package model

import (
	"database/sql/driver"
	"strings"
	"time"
)

type Strs []string

func (m *Strs) Scan(val interface{}) error {
	s := val.([]uint8)
	ss := strings.Split(string(s), "|")
	*m = ss
	return nil
}

func (m *Strs) Value() (driver.Value, error) {
	str := strings.Join(*m, "|")
	return str, nil
}

type Field struct {
	EventID     string    `gorm:"primaryKey" json:"eventID" binding:"required,max=32"`
	ID          string    `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Type        string    `json:"type"`
	TypeID      string    `json:"typeID"`
	Key         string    `json:"key"`
	Value       Strs      `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
