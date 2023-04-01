package model

import "time"

type FieldLog struct {
	EventLogID string `gorm:"index;not null" json:"eventLogID" binding:"required,max=32"`
	FieldID    string `gorm:"index;not null" json:"fieldID" binding:"required,max=32"`
	ID         string `gorm:"type:varchar(191);primaryKey" json:"id" binding:"required,max=32"`

	Key   string `json:"key"`
	Value string `json:"value"`

	Used     bool   `gorm:"default:false" json:"used"`   //被使用过（被前端clear log）
	Tested   bool   `gorm:"default:false" json:"tested"` //被测试用
	Platform string `gorm:"not null" json:"platform"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
