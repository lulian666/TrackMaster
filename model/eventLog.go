package model

import "time"

type EventLog struct {
	EventID string `gorm:"index;not null" json:"eventID" binding:"required,max=32"`
	ID      string `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Name    string `gorm:"not null" json:"name" binding:"required,min=2,max=50"`

	Used     bool   `gorm:"default:false" json:"used"`   //被使用过（被前端clear log）
	Tested   bool   `gorm:"default:false" json:"tested"` //被测试用
	Platform string `gorm:"not null" json:"platform"`
	Raw      string `gorm:"not null;size:2000;type:json" json:"raw"`
	Content  string `gorm:"not null;size:200;type:json" json:"content"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	FieldLogs []FieldLog `gorm:"foreignKey:EventLogID"`
}
