package service

import (
	"TrackMaster/model"
	"TrackMaster/third_party/jet"
	"gorm.io/gorm"
)

type ProjectService interface {
	SyncProject() error
	ListProject() ([]model.Project, error)
}

type projectService struct {
	DB *gorm.DB
}

func NewProjectService(db *gorm.DB) ProjectService {
	return &projectService{
		DB: db,
	}
}

// SyncProject
// 获取project list
// 和本地本地的对比
// 只新增不存在的id
func (s projectService) SyncProject() error {
	projects, err := jet.GetProjects()
	if err != nil {
		return err
	}

	projectIDs := make([]string, len(projects))

	for i, project := range projects {
		projectIDs[i] = project.ID
	}

	var existingProjects []model.Project
	s.DB.Where("id IN ?", projectIDs).Find(&existingProjects)

	for _, project := range projects {
		var p model.Project
		for _, existingProject := range existingProjects {
			if existingProject.ID == project.ID {
				p = existingProject
			}
		}

		if p.ID == "" {
			p.ID = project.ID
			p.Name = project.CnName
			err := p.Create(s.DB)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s projectService) ListProject() ([]model.Project, error) {
	p := &model.Project{}
	return p.List(s.DB)
}
