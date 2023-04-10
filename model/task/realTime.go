package task

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/pkg/track"
	"TrackMaster/pkg/worker"
	"fmt"
)

type Duty string

const (
	Fetch Duty = "Fetch"
	Test  Duty = "Test"
)

type RealTimeTrackTask struct {
	Type   Duty
	Record model.Record
	WP     *worker.Pool
}

func (t *RealTimeTrackTask) Execute() *pkg.Error {
	if t.Type == Fetch {
		count, err := track.FetchNewLog(t.Record)
		if err != nil {
			return err
		}

		if count > 0 {
			testJob := &RealTimeTrackTask{
				Type:   Test,
				Record: t.Record,
			}
			select {
			case t.WP.Jobs <- testJob:
				fmt.Println("Test job has been added to the queue")
			default:
				fmt.Println("Too many jobs")
			}
		}
	}

	if t.Type == Test {
		err := track.TestNewLog(t.Record)
		if err != nil {
			return err
		}
	}

	return nil
}
