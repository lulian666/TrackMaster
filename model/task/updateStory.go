package task

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/pkg/track"
	"TrackMaster/pkg/worker"
	"time"
)

type UpdateStoryTask struct {
	Project  *model.Project
	Interval time.Duration
	WP       *worker.Pool
}

func (t *UpdateStoryTask) Execute() *pkg.Error {
	stories, err := track.UpdateStory(t.Project)
	if err != nil {
		return err
	}

	// fix: story需要先保存
	// 更新events
	for i := range stories {
		go func(i int) {
			err := track.SyncEvent(stories[i].ID)
			if err != nil {
				t.WP.Errors <- err
			}
		}(i)
	}
	return nil
}
