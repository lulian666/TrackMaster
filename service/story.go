package service

import (
	"TrackMaster/model"
	"gorm.io/gorm"
)

type StoryService interface {
	SyncStory(p *model.Project) error
}

type storyService struct {
	DB *gorm.DB
}

func NewStoryService(db *gorm.DB) StoryService {
	return &storyService{
		DB: db,
	}
}

func (s storyService) SyncStory(p *model.Project) error {
	// project 是否存在
	err := p.Get(s.DB)
	if err != nil {
		return err
	}
	// 先同步enum type，再同步story
	enumTypeS := NewEnumTypeService(s.DB)
	err = enumTypeS.SyncEnumType(p)
	if err != nil {
		return err
	}
	return nil
}
