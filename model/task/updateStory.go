package task

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/pkg/track"
	"TrackMaster/pkg/worker"
	"TrackMaster/third_party/jet"
	"time"
)

type UpdateStoryTask struct {
	Project  *model.Project
	Interval time.Duration
	WP       *worker.Pool
}

func (t *UpdateStoryTask) Execute() *pkg.Error {
	// 先同步enum type，再同步story
	err := track.SyncEnumType(t.Project)
	if err != nil {
		return err
	}

	stories, err := jet.GetStories(t.Project.ID)
	if err != nil {
		return err
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
				st.ProjectID = t.Project.ID
				st.ID = stories[i].ID
				if stories[i].Name != existingStories[j].Name || stories[i].Desc != existingStories[j].Description {
					st.Name = stories[i].Name
					st.Description = stories[i].Desc
					storyUpdateList = append(storyUpdateList, st)
				}
			}
		}

		if st.ID == "" {
			st.ProjectID = t.Project.ID
			st.ID = stories[i].ID
			st.Name = stories[i].Name
			st.Description = stories[i].Desc
			storyCreateList = append(storyCreateList, st)
		}

	}

	if len(storyUpdateList) > 0 {
		result := initializer.DB.Save(storyUpdateList)
		if result.Error != nil {
			return err
		}
	}

	if len(storyCreateList) > 0 {
		result := initializer.DB.Create(storyCreateList)
		if result.Error != nil {
			return err
		}
	}

	// fix: story需要先保存
	// 更新events
	for i := range stories {
		go func(i int) {
			err := SyncEventHardcore(stories[i].ID)
			if err != nil {
				t.WP.Errors <- err
			}
		}(i)
	}

	return nil
}

func SyncEventHardcore(storyID string) *pkg.Error {
	events, err := jet.GetEvents(storyID)
	if err != nil {
		return err
	}

	for i := range events {
		// sync fields hardcore
		fields := events[i].EventDefinitions

		var existingFields []model.Field
		fieldIDs := make([]string, len(fields))
		for k := range fields {
			fieldIDs[k] = fields[k].ID
		}
		// 先删除子表中关联的记录
		var results []model.FieldResult
		initializer.DB.Delete(&results, "field_id IN (?)", fieldIDs)
		initializer.DB.Delete(&existingFields, "event_id = ?", events[i].ID)
	}

	var existingEvents []model.Event
	// delete all events that match story id
	result := initializer.DB.Delete(&existingEvents, "story_id = ?", storyID)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}

	eventCreateList := make([]model.Event, 0, len(events))
	fieldCreateList := make([]model.Field, 0, len(events)*6)

	for i := range events {
		e := model.Event{
			StoryID:     storyID,
			ID:          events[i].ID,
			Name:        events[i].Name,
			Description: events[i].Desc,
		}
		eventCreateList = append(eventCreateList, e)
		fields := events[i].EventDefinitions

		for m := range fields {
			f := model.Field{}

			//需要创建的记录
			f.EventID = e.ID
			f.ID = fields[m].ID
			f.Type = fields[m].Type.Name
			f.TypeID = fields[m].Type.ID
			f.Key = fields[m].Name
			value, err := track.LocateValue(fields[m], initializer.DB)
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

	// 批量创建
	if len(eventCreateList) > 0 {
		result := initializer.DB.Create(eventCreateList)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}
	}

	if len(fieldCreateList) > 0 {
		result := initializer.DB.Create(fieldCreateList)
		if result.Error != nil {
			return pkg.NewError(pkg.ServerError, result.Error.Error())
		}
	}
	return nil
}
