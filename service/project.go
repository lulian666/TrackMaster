package service

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/third_party/jet"
	"gorm.io/gorm"
)

type ProjectService interface {
	SyncProject() *pkg.Error
	ListProject() ([]model.Project, *pkg.Error)
}

type projectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) ProjectService {
	return &projectService{
		db: db,
	}
}

// SyncProject
// 获取project list
// 和本地本地的对比
// 只新增不存在的id
func (s projectService) SyncProject() *pkg.Error {
	projects, err := jet.GetProjects()
	if err != nil {
		return err
	}

	projectIDs := make([]string, len(projects))

	for i, project := range projects {
		projectIDs[i] = project.ID
	}

	var existingProjects []model.Project
	s.db.Where("id IN ?", projectIDs).Find(&existingProjects)

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
			err := p.Create(s.db)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s projectService) ListProject() ([]model.Project, *pkg.Error) {
	p := &model.Project{}
	return p.List(s.db)
}
