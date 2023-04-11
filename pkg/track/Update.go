package track

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/third_party/jet"
	"reflect"
)

func UpdateStory(project *model.Project) ([]jet.Story, *pkg.Error) {
	// 先同步enum type，再同步story
	err := SyncEnumType(project)
	if err != nil {
		return nil, err
	}

	stories, err := jet.GetStories(project.ID)
	if err != nil {
		return nil, err
	}

	storyIDs := make([]string, len(stories))
	for i, st := range stories {
		storyIDs[i] = st.ID
	}

	var existingStories []model.Story
	initializer.DB.Where("id in (?)", storyIDs).Find(&existingStories)

	storyCreateList := make([]model.Story, 0, len(stories))
	storyUpdateList := make([]model.Story, 0, len(stories))
	for i := range stories {
		st := model.Story{}
		for j := range existingStories {
			if stories[i].ID == existingStories[j].ID {
				st.ProjectID = project.ID
				st.ID = stories[i].ID
				if stories[i].Name != existingStories[j].Name || stories[i].Desc != existingStories[j].Description {
					st.Name = stories[i].Name
					st.Description = stories[i].Desc
					storyUpdateList = append(storyUpdateList, st)
				}
			}
		}

		if st.ID == "" {
			st.ProjectID = project.ID
			st.ID = stories[i].ID
			st.Name = stories[i].Name
			st.Description = stories[i].Desc
			storyCreateList = append(storyCreateList, st)
		}

	}

	if len(storyUpdateList) > 0 {
		result := initializer.DB.Save(storyUpdateList)
		if result.Error != nil {
			return nil, err
		}
	}

	if len(storyCreateList) > 0 {
		result := initializer.DB.Create(storyCreateList)
		if result.Error != nil {
			return nil, err
		}
	}

	return stories, nil
}

func SyncEvent(storyID string) *pkg.Error {
	events, err := jet.GetEvents(storyID)
	if err != nil {
		return err
	}

	// 删除需求中已经不存在的events和fields
	err = deleteNonExistingEvents(storyID, events)
	if err != nil {
		return err
	}

	eventIDs := make([]string, len(events))
	for i := range events {
		eventIDs[i] = events[i].ID
	}

	var existingEvents []model.Event
	initializer.DB.Where("id IN (?)", eventIDs).Find(&existingEvents)

	eventCreateList := make([]model.Event, 0, len(events))
	eventUpdateList := make([]model.Event, 0, len(events))

	fieldCreateList := make([]model.Field, 0, len(events)*6)
	fieldUpdateList := make([]model.Field, 0, len(events)*6)

	for i := range events {
		e := model.Event{}
		for j := range existingEvents {
			if events[i].ID == existingEvents[j].ID && existingEvents[j].StoryID == storyID {
				e.StoryID = storyID
				e.ID = events[i].ID
				if events[i].Name != existingEvents[j].Name || events[i].Desc != existingEvents[j].Description {
					e.Name = events[i].Name
					e.Description = events[i].Desc
					eventUpdateList = append(eventUpdateList, e)
				}
			}
		}

		if e.ID == "" {
			e.StoryID = storyID
			e.ID = events[i].ID
			e.Name = events[i].Name
			e.Description = events[i].Desc
			eventCreateList = append(eventCreateList, e)
		}

		// sync fields
		fields := events[i].EventDefinitions

		var existingFields []model.Field
		initializer.DB.Where("event_id IN (?)", eventIDs).Find(&existingFields)

		for m := range fields {
			f := model.Field{}
			for n := range existingFields {
				// 不仅field的id要一样， event id也要一样
				if fields[m].ID == existingFields[n].ID && existingFields[n].EventID == e.ID {
					f.EventID = e.ID
					f.ID = fields[m].ID
					if anyDifference(existingFields[n], fields[m]) {
						value, err := LocateValue(fields[m], initializer.DB)
						if err != nil {
							return err
						}
						f.Value, err = pkg.Strs(value).Value()
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
				value, err := LocateValue(fields[m], initializer.DB)
				if err != nil {
					return err
				}
				f.Value, err = pkg.Strs(value).Value()
				if err != nil {
					return err
				}
				f.Description = fields[m].Note
				fieldCreateList = append(fieldCreateList, f)
			}
		}

	}
	if len(eventUpdateList) > 0 {
		result := initializer.DB.Save(eventUpdateList)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}
	}

	if len(eventCreateList) > 0 {
		result := initializer.DB.Create(eventCreateList)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}
	}

	// 批量更新
	if len(fieldUpdateList) > 0 {
		result := initializer.DB.Save(fieldUpdateList)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}
	}

	// 批量创建
	if len(fieldCreateList) > 0 {
		result := initializer.DB.Create(fieldCreateList)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}
	}
	return nil

}

func deleteNonExistingEvents(storyID string, events []jet.Event) *pkg.Error {
	var existingEventIDs []string
	result := initializer.DB.Model(&model.Event{}).Where("story_id = ?", storyID).Pluck("id", &existingEventIDs)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}

	newEventIDs := make([]string, len(events))
	for i, event := range events {
		newEventIDs[i] = event.ID
	}

	deletedEventIDs := sliceDiff(existingEventIDs, newEventIDs)
	if len(deletedEventIDs) > 0 {
		// 不在的先删掉字表里的记录，再删掉event
		var deletedFieldIDs []string
		result = initializer.DB.Model(&model.Field{}).Where("event_id IN (?)", deletedEventIDs).Pluck("id", &deletedFieldIDs)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}

		result = initializer.DB.Delete(&model.FieldResult{}, "field_id IN (?)", deletedFieldIDs)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}

		result = initializer.DB.Delete(&model.EventResult{}, "event_id IN (?)", deletedEventIDs)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}

		result = initializer.DB.Exec("DELETE FROM records WHERE id IN (SELECT record_id FROM record_events WHERE event_id IN (?))", deletedEventIDs)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}

		result = initializer.DB.Delete(&model.Field{}, "id in (?)", deletedFieldIDs)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}

		result = initializer.DB.Delete(&model.Event{}, "id in (?)", deletedEventIDs)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}
	}

	errCh := make(chan *pkg.Error, 1)
	for i := range events {
		go func(i int) {
			err := deleteNonExistingFields(events[i].ID, events[i].EventDefinitions)
			if err != nil {
				errCh <- err
			}
		}(i)
	}

	select {
	case err := <-errCh:
		return err
	default:
		// do nothing
	}
	return nil
}

func deleteNonExistingFields(eventID string, fields []jet.Field) *pkg.Error {
	var existingFields []model.Field
	initializer.DB.Model(&model.Field{}).Where("event_id = ?", eventID).Select("id", "event_id").Find(&existingFields)

	newFields := make([]model.Field, len(fields))
	for i, field := range fields {
		newFields[i] = model.Field{
			ID:      field.ID,
			EventID: eventID,
		}
	}

	deletedFields := findMissingItems(existingFields, newFields)
	if len(deletedFields) > 0 {
		deletedFieldIDs := make([]string, len(fields))
		for i, field := range deletedFields {
			deletedFieldIDs[i] = field.ID
		}

		result := initializer.DB.Delete(&model.FieldResult{}, "field_id IN (?)", deletedFieldIDs)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}

		result = initializer.DB.Delete(&model.Field{}, "id in (?)", deletedFieldIDs)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}
	}

	return nil
}

func sliceDiff(source []string, target []string) []string {
	var diff []string
	for _, t := range target {
		found := false
		for _, s := range source {
			if t == s {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, t)
		}
	}
	return diff
}

func findMissingItems(s1 []model.Field, s2 []model.Field) []model.Field {
	var missingItems []model.Field

	for _, item1 := range s1 {
		found := false
		for _, item2 := range s2 {
			if item1.ID == item2.ID && item1.EventID == item2.EventID {
				found = true
				break
			}
		}
		if !found {
			missingItems = append(missingItems, item1)
		}
	}

	return missingItems
}

func anyDifference(existingF model.Field, f jet.Field) bool {
	if existingF.Type != f.Type.Name ||
		existingF.TypeID != f.Type.ID ||
		existingF.Description != f.Note ||
		!reflect.DeepEqual(existingF.Value, f.Values) {
		return true
	}
	return false
}
