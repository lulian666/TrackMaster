package service

import (
	"TrackMaster/model"
	"TrackMaster/third_party/jet"
	"gorm.io/gorm"
)

type EventService interface {
	SyncEvent(story model.Story) error
}

type eventService struct {
	db *gorm.DB
}

func NewEventService(db *gorm.DB) EventService {
	return &eventService{db: db}
}

func (s eventService) SyncEvent(story model.Story) error {
	events, err := jet.GetEvents(story.ID)
	if err != nil {
		return err
	}

	eventIDs := make([]string, len(events))
	for i := range events {
		eventIDs[i] = events[i].ID
	}

	var existingEvents []model.Event
	s.db.Where("id IN (?)", eventIDs).Find(&existingEvents)

	eventCreateList := make([]model.Event, 0, len(events))
	eventUpdateList := make([]model.Event, 0, len(events))

	for i := range events {
		e := model.Event{}
		for j := range existingEvents {
			if events[i].ID == existingEvents[j].ID && existingEvents[j].StoryID == story.ID {
				e.StoryID = story.ID
				e.ID = events[i].ID
				if events[i].Name != existingEvents[j].Name || events[i].Desc != existingEvents[j].Description {
					e.Name = events[i].Name
					e.Description = events[i].Desc
					eventUpdateList = append(eventUpdateList, e)
				}
			}
		}

		if e.ID == "" {
			e.StoryID = story.ID
			e.ID = events[i].ID
			e.Name = events[i].Name
			e.Description = events[i].Desc
			eventCreateList = append(eventCreateList, e)
		}

		// sync fields
		fields := events[i].EventDefinitions

		var existingFields []model.Field
		s.db.Where("event_id IN (?)", eventIDs).Find(&existingFields)

		fieldCreateList := make([]model.Field, 0, len(fields))
		fieldUpdateList := make([]model.Field, 0, len(fields))

		for m := range fields {
			f := model.Field{}
			for n := range existingFields {
				// 不仅field的id要一样， event id也要一样
				if fields[m].ID == existingFields[n].ID && existingFields[n].EventID == e.ID {
					f.EventID = e.ID
					f.ID = fields[m].ID
					if anyDifference(existingFields[n], fields[m]) {
						f.Value, err = locateValue(fields[m], s.db)
						if err != nil {
							return err
						}
						f.Type = fields[m].Type.Name
						f.TypeID = fields[m].Type.ID
						f.Key = fields[m].Name
						f.Description = fields[m].Note
						fieldUpdateList = append(fieldUpdateList, f)
					}
				}
			}
			//需要创建的记录
			if f.ID == "" {
				f.EventID = e.ID
				f.ID = fields[m].ID
				f.Type = fields[m].Type.Name
				f.TypeID = fields[m].Type.ID
				f.Key = fields[m].Name
				f.Value, err = locateValue(fields[m], s.db)
				if err != nil {
					return err
				}
				f.Description = fields[m].Note
				fieldCreateList = append(fieldCreateList, f)
			}
		}
		// 批量更新
		if len(fieldUpdateList) > 0 {
			result := s.db.Save(fieldUpdateList)
			if result.Error != nil {
				return result.Error
			}
		}

		// 批量创建
		if len(fieldCreateList) > 0 {
			result := s.db.Create(fieldCreateList)
			if result.Error != nil {
				return result.Error
			}
		}
	}
	if len(eventUpdateList) > 0 {
		result := s.db.Save(eventUpdateList)
		if result.Error != nil {
			return result.Error
		}
	}

	if len(eventCreateList) > 0 {
		result := s.db.Create(eventCreateList)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func anyDifference(existingF model.Field, f jet.Field) bool {
	if existingF.Type != f.Type.Name ||
		existingF.TypeID != f.Type.ID ||
		existingF.Description != f.Note {
		return true
	}
	return false
}

func locateValue(field jet.Field, db *gorm.DB) ([]string, error) {
	if len(field.Values) > 0 {
		// 根据type id和id去拿值
		values := make([]string, 0, len(field.Values))
		for _, v := range values {
			value := model.EnumValue{}
			result := db.Where("type_id = ?", field.Type.ID).Where("id = ?", v).Find(&value)
			if result.Error != nil {
				return nil, result.Error
			}
			if value.ID != "" {
				values = append(values, value.Name)
			}
		}
		return values, nil
	} else {
		// 没有值
		return []string{""}, nil
	}
}
