package model

import (
	"TrackMaster/pkg"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	StoryID     string        `gorm:"index;not null" json:"storyID" binding:"required,max=32"`
	ID          string        `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Name        string        `gorm:"not null" json:"name" binding:"required,min=2,max=50"`
	Description string        `json:"description"`
	OnTrail     bool          `gorm:"default:false" json:"onTrail"` //这个字段暂时用不上
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updatedAt"`
	Fields      []Field       `gorm:"foreignKey:EventID"`
	Results     []EventResult `gorm:"foreignKey:EventID"`

	RecordIDs []string  `gorm:"-"` // 不要在数据库中保存这个字段
	Records   []*Record `gorm:"many2many:record_events;"`
}

type SwaggerEvents struct {
	Data  []*Event
	Pager *pkg.Pager
}

func (e *Event) List(db *gorm.DB, eventIDs []string) ([]Event, int64, *pkg.Error) {
	var events []Event
	result := db.Model(Event{}).Preload("Fields").Where("id in (?)", eventIDs).Find(&events)

	if result.Error != nil {
		return nil, 0, pkg.NewError(pkg.ServerError, result.Error.Error())
	}

	totalRow := result.RowsAffected
	return events, totalRow, nil
}

func (e *Event) ListWithNewestResult(db *gorm.DB, eventIDs []string, recordID string) ([]Event, int64, *pkg.Error) {
	var events []Event
	result := db.Model(e).Preload("Fields.Results", func(db *gorm.DB) *gorm.DB {
		return db.Where("record_id = ?", recordID).Order("created_at desc")
	}).Preload("Results", func(db *gorm.DB) *gorm.DB {
		return db.Where("record_id = ?", recordID).Order("created_at desc")
	}).Where("id in (?)", eventIDs).Find(&events)

	totalRow := result.RowsAffected
	return events, totalRow, nil
}

func (e *Event) ListEventName(db *gorm.DB, eventIDs []string) ([]string, int64, *pkg.Error) {
	var events []string
	result := db.Model(e).Where("id in (?)", eventIDs).Pluck("name", &events)

	if result.Error != nil {
		return nil, 0, pkg.NewError(pkg.ServerError, result.Error.Error())
	}

	totalRow := result.RowsAffected
	return events, totalRow, nil
}

type Events []Event

func (es Events) FindByID(id string) (Event, bool) {
	for i := range es {
		if es[i].ID == id {
			return es[i], true
		}
	}
	return Event{}, false
}

func (es Events) ListEventID() ([]string, *pkg.Error) {
	eventIDs := make([]string, 0, len(es))
	for i := range es {
		eventIDs = append(eventIDs, es[i].ID)
	}
	return eventIDs, nil
}
