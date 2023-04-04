package model

import (
	"time"
)

type Field struct {
	EventID     string    `gorm:"type:varchar(191);primaryKey" json:"eventID" binding:"required,max=32"`
	ID          string    `gorm:"type:varchar(191);primaryKey" json:"id" binding:"required,max=32"`
	Type        string    `json:"type"`
	TypeID      string    `json:"typeID"`
	Key         string    `json:"key"`
	Value       string    `json:"value"` // 插入时会把数组用"|"隔开变成字符串，读取的时候也要转化一下
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Fields []Field

func (fs Fields) FindByID(id string) (Field, bool) {
	for i := range fs {
		if fs[i].ID == id {
			return fs[i], true
		}
	}
	return Field{}, false
}
