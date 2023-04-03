package service

import "gorm.io/gorm"

type EventLogService interface {
	InComingLog()
}

type eventLogService struct {
	db *gorm.DB
}

func NewEventLogService(db *gorm.DB) EventService {
	return &eventService{
		db: db,
	}
}

func (s eventLogService) InComingLog() {

}
