package model

import (
	"gorm.io/gorm"
	"time"
)

type EventLog struct {
	RecordID string `gorm:"index;not null" json:"recordID" binding:"required,max=32"`
	EventID  string `gorm:"index;not null" json:"eventID" binding:"required,max=32"`
	ID       string `gorm:"primaryKey" json:"id" binding:"required,max=32"` // jet返回时每个eventLog都带id
	Name     string `gorm:"not null" json:"name" binding:"required,min=2,max=50"`
	UserID   string `json:"userId"` // 既然有这个信息，不妨存一下

	Used     bool   `gorm:"default:false" json:"used"`   //被使用过（被前端clear log）
	Tested   bool   `gorm:"default:false" json:"tested"` //被测试用
	Platform string `gorm:"not null" json:"platform"`
	Raw      string `gorm:"not null;size:2000;type:json" json:"raw"`
	Content  string `gorm:"size:200;type:json;default:null" json:"content"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	FieldLogs []FieldLog `gorm:"foreignKey:EventLogID"`
}

func (e *EventLog) ListUnused(db *gorm.DB, recordID string) ([]EventLog, int64, error) {
	var eventLogs []EventLog
	result := db.Preload("FieldLogs").Where("record_id = ?", recordID).Where("used = ?", false).Order("created_at desc").Find(&eventLogs)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	totalRow := result.RowsAffected
	return eventLogs, totalRow, nil
}

type EventLogs []EventLog

func (e EventLogs) UpdateToUsed(db *gorm.DB) error {
	result := db.Save(e).Update("used", true)
	return result.Error
}
