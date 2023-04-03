package service

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/third_party/jet"
	"gorm.io/gorm"
)

type StoryService interface {
	SyncStory(p *model.Project) error
	ListStory(story *model.Story, pager pkg.Pager) ([]model.Story, int64, error)
	GetStory(story *model.Story) error
}

type storyService struct {
	db *gorm.DB
}

func NewStoryService(db *gorm.DB) StoryService {
	return &storyService{
		db: db,
	}
}

func (s storyService) SyncStory(p *model.Project) error {
	// project 是否存在
	err := p.Get(s.db)
	if err != nil {
		return err
	}
	// 先同步enum type，再同步story
	enumTypeS := NewEnumTypeService(s.db)
	err = enumTypeS.SyncEnumType(p)
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
	storyList := append(storyCreateList, storyUpdateList...)
	for i := range storyList {
		err = eventS.SyncEvent(storyList[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s storyService) ListStory(story *model.Story, pager pkg.Pager) ([]model.Story, int64, error) {
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

func (s storyService) GetStory(story *model.Story) error {
	return story.Get(s.db)
}
