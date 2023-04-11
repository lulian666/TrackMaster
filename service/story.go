package service

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/pkg/track"
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

	stories, err := track.UpdateStory(p)
	if err != nil {
		return err
	}

	// fix: story需要先保存
	// 更新events
	errCh := make(chan *pkg.Error, 1)
	for i := range stories {
		go func(i int) {
			err := track.SyncEvent(stories[i].ID)
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
