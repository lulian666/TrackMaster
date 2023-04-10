package service

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/pkg/track"
	"TrackMaster/third_party/jet"
	"gorm.io/gorm"
)

type StoryService interface {
	SyncStory(p *model.Project) *pkg.Error
	ListStory(story *model.Story, pager pkg.Pager) ([]model.Story, int64, *pkg.Error)
	GetStory(story *model.Story) *pkg.Error
}

type storyService struct {
	db *gorm.DB
}

func NewStoryService(db *gorm.DB) StoryService {
	return &storyService{
		db: db,
	}
}

func (s storyService) SyncStory(p *model.Project) *pkg.Error {
	// project 是否存在
	err := p.Get(s.db)
	if err != nil {
		return err
	}
	// 先同步enum type，再同步story
	err = track.SyncEnumType(p)
	if err != nil {
		return err
	}

	stories, err := jet.GetStories(p.ID)
	if err != nil {
		return err
	}

	storyIDs := make([]string, len(stories))
	for i, st := range stories {
		storyIDs[i] = st.ID
	}

	var existingStories []model.Story
	s.db.Where("id in (?)", storyIDs).Find(&existingStories)

	storyCreateList := make([]model.Story, 0, len(stories))
	storyUpdateList := make([]model.Story, 0, len(stories))
	for i := range stories {
		st := model.Story{}
		for j := range existingStories {
			if stories[i].ID == existingStories[j].ID {
				st.ProjectID = p.ID
				st.ID = stories[i].ID
				if stories[i].Name != existingStories[j].Name || stories[i].Desc != existingStories[j].Description {
					st.Name = stories[i].Name
					st.Description = stories[i].Desc
					storyUpdateList = append(storyUpdateList, st)
				}
			}
		}

		if st.ID == "" {
			st.ProjectID = p.ID
			st.ID = stories[i].ID
			st.Name = stories[i].Name
			st.Description = stories[i].Desc
			storyCreateList = append(storyCreateList, st)
		}

	}

	if len(storyUpdateList) > 0 {
		result := s.db.Save(storyUpdateList)
		if result.Error != nil {
			return err
		}
	}

	if len(storyCreateList) > 0 {
		result := s.db.Create(storyCreateList)
		if result.Error != nil {
			return err
		}
	}

	// fix: story需要先保存
	// 更新events
	eventS := NewEventService(s.db)
	errCh := make(chan *pkg.Error, 1)
	for i := range stories {
		go func(i int) {
			err := eventS.SyncEvent(stories[i].ID)
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

func (s storyService) ListStory(story *model.Story, pager pkg.Pager) ([]model.Story, int64, *pkg.Error) {
	// 先判断project是否存在
	project := model.Project{
		ID: story.ProjectID,
	}
	err := project.Get(s.db)
	if err != nil {
		return nil, 0, err
	}

	pageOffset := pkg.GetPageOffset(pager.Page, pager.PageSize)
	return story.List(s.db, pageOffset, pager.PageSize)
}

func (s storyService) GetStory(story *model.Story) *pkg.Error {
	return story.Get(s.db)
}
