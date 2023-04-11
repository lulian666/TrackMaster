package service

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"TrackMaster/model/task"
	"TrackMaster/pkg"
	"TrackMaster/pkg/worker"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ScheduleService interface {
	On(p *model.Project, wp *worker.Pool) (model.Schedule, *pkg.Error)
	Off(p *model.Project) (model.Schedule, *pkg.Error)
}

type scheduleService struct {
	db *gorm.DB
}

func NewScheduleService(db *gorm.DB) ScheduleService {
	return &scheduleService{
		db: db,
	}
}

func (s scheduleService) On(p *model.Project, wp *worker.Pool) (model.Schedule, *pkg.Error) {
	err := p.Get(s.db)
	if err != nil {
		if strings.Contains(err.Msg, "record not found") {
			return model.Schedule{}, pkg.NewError(pkg.BadRequest, "project does not exist")
		}
		return model.Schedule{}, err
	}

	schedule := model.Schedule{
		ProjectID: p.ID,
	}
	err = schedule.Get(s.db)
	if err != nil {
		if strings.Contains(err.Msg, "record not found") {
			schedule := model.Schedule{
				ProjectID: p.ID,
				Interval:  3 * time.Hour,
				Name:      p.Name + "_定期更新需求",
				Status:    true,
			}

			err = schedule.Create(s.db)
			if err != nil {
				return schedule, err
			}

			// 如果不存在则开启一个新任务
			limit := 3000
			go updateStory(limit, &schedule, p, wp)
		}
		return schedule, err
	}

	// 如果本身已经开启，则不再开启新的任务
	if !schedule.Status {
		schedule.Status = true
		s.db.Save(schedule)

		limit := 3000
		go updateStory(limit, &schedule, p, wp)
	}

	return schedule, nil
}

func (s scheduleService) Off(p *model.Project) (model.Schedule, *pkg.Error) {
	err := p.Get(s.db)
	if err != nil {
		if strings.Contains(err.Msg, "record not found") {
			return model.Schedule{}, pkg.NewError(pkg.BadRequest, "project does not exist")
		}
		return model.Schedule{}, err
	}

	schedule := model.Schedule{
		ProjectID: p.ID,
	}
	err = schedule.Get(s.db)
	if err != nil {
		return model.Schedule{}, err
	}

	// 如果本身已经关闭，则不进行修改
	if schedule.Status {
		schedule.Status = false
		s.db.Save(schedule)

	}
	return schedule, nil
}

func updateStory(limit int, schedule *model.Schedule, p *model.Project, wp *worker.Pool) {
	d := 3 * time.Hour
	// 保护
	if schedule.Interval >= 1*time.Hour && schedule.Interval <= 24*time.Hour {
		d = schedule.Interval
	}

	ticker := time.NewTicker(d)
	i := 1

	for range ticker.C {
		initializer.DB.First(schedule)

		if i > limit || schedule.Status == false {
			result := initializer.DB.Model(schedule).Update("status", false)
			if result.Error != nil {
				return
			}

			ticker.Stop()
			break
		}

		job := task.UpdateStoryTask{
			Project:  p,
			Interval: schedule.Interval,
			WP:       wp,
		}

		select {
		case wp.Jobs <- &job:
			fmt.Println("job has been added to the queue")
		default:
			fmt.Println("Too many jobs")
		}
		i += 1
	}
}
