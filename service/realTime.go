package service

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"TrackMaster/model/request"
	"TrackMaster/model/task"
	"TrackMaster/pkg"
	"TrackMaster/pkg/worker"
	"TrackMaster/third_party/jet"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
	"strings"
	"time"
)

type RealTimeService interface {
	Start(wp *worker.Pool, req request.Start) (model.Record, *pkg.Error)
	Stop(r *model.Record) *pkg.Error
	Update(r *model.Record, req request.Update) *pkg.Error
	GetLog(r *model.Record) ([]model.EventLog, int64, *pkg.Error)
	ClearLog(r *model.Record) *pkg.Error
	UpdateResult(req request.UpdateResult) *pkg.Error
	GetResult(r *model.Record) ([]model.Event, int64, *pkg.Error)
	Test(wp *worker.Pool, req request.Start) (model.Record, *pkg.Error)
	//recordExist(r *model.Record) *pkg.Error
}

type realTimeService struct {
	db *gorm.DB
}

func NewRealTimeService(db *gorm.DB) RealTimeService {
	return &realTimeService{
		db: db,
	}
}

func (s realTimeService) Start(wp *worker.Pool, req request.Start) (model.Record, *pkg.Error) {
	// 判断events合不合法
	_, eventNames, err, events := eventsLegitimate(s.db, req.EventIDs)
	if err != nil {
		return model.Record{}, err
	}
	// 判断accounts合不合法
	err = accountsLegitimate(s.db, req.AccountIDs)
	if err != nil {
		return model.Record{}, err
	}

	// 判断project存不存在
	p := model.Project{
		ID: req.Project,
	}
	err = p.Get(s.db)
	if err != nil {
		if strings.Contains(err.Msg, "record not found") {
			return model.Record{}, pkg.NewError(pkg.BadRequest, "project does not exist")
		}
		return model.Record{}, err
	}

	// 创建filter
	filter := jet.Filter{
		Events:  eventNames,
		Project: p.ID,
		UserIDs: req.AccountIDs,
	}
	filterRes, err := filter.Create()
	if err != nil {
		return model.Record{}, err
	}

	if filterRes.Status != "READY" {
		return model.Record{}, pkg.NewError(pkg.ServerError, "调用jet创建filter时出了问题，请稍后再试")
	}

	filter.ID = filterRes.ID
	filter.Status = jet.RECORDING
	err = filter.Update()
	if err != nil {
		return model.Record{}, err
	}

	// 创建record
	u := uuid.New()
	id := strings.ReplaceAll(u.String(), "-", "")

	r := model.Record{
		Name:      "实时埋点测试" + time.Stamp,
		Status:    model.ON,
		Filter:    filter.ID,
		ProjectID: filter.Project,
		ID:        id,
		//Events:    eventsValue,
		Events: events,
	}
	err = r.Create(s.db)
	if err != nil {
		return r, err
	}

	// todo 两种方式
	// 1.先尝试用go的标准库time和channel来实现试试
	// 2.再尝试加上全局的workerPool
	//logCh := make(chan int)
	l := os.Getenv("REALTIME_LIMIT")
	limit := pkg.StrTo(l).MustInt()

	go checkLog(limit, r, wp)
	//go testLog(logCh, r, wp)

	return r, nil
}

func (s realTimeService) Stop(r *model.Record) *pkg.Error {
	err := r.Update(s.db)
	if err != nil {
		if strings.Contains(err.Msg, "record not found") {
			return pkg.NewError(pkg.NotFound, fmt.Sprintf("record with id %s not found", r.ID))
		}
		return err
	}

	filter := jet.Filter{
		ID:     r.Filter,
		Status: jet.STOPPED,
	}
	err = filter.Update()
	if err != nil {
		return err
	}

	return nil
}

func (s realTimeService) Update(r *model.Record, req request.Update) *pkg.Error {
	// 判断record是否存在
	err := s.recordExist(r)
	if err != nil {
		return err
	}

	// 修改filter
	filter := jet.Filter{
		ID: r.Filter,
	}

	// 判断events合不合法
	if len(req.EventIDs) > 0 {
		_, eventNames, err, events := eventsLegitimate(s.db, req.EventIDs)
		if err != nil {
			return err
		}
		filter.Events = eventNames
		//eventsValue, _ := pkg.Strs(eventIDs).Value()
		//r.Events = eventsValue
		r.Events = events
	}

	// 判断accounts合不合法
	if len(req.AccountIDs) > 0 {
		err := accountsLegitimate(s.db, req.AccountIDs)
		if err != nil {
			return err
		}
		filter.UserIDs = req.AccountIDs
	}

	err = filter.Update()
	if err != nil {
		return err
	}

	result := s.db.Save(&r)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}

func (s realTimeService) GetLog(r *model.Record) ([]model.EventLog, int64, *pkg.Error) {
	// 只取未used的log并且按创建时间倒叙排
	el := model.EventLog{}
	logs, totalRow, err := el.ListUnused(s.db, r.ID)

	if err != nil {
		if strings.Contains(err.Msg, "record not found") {
			return nil, 0, pkg.NewError(pkg.NotFound, fmt.Sprintf("record with id %s not found", r.ID))
		}
		return nil, 0, err
	}

	return logs, totalRow, nil
}

func (s realTimeService) ClearLog(r *model.Record) *pkg.Error {
	err := s.recordExist(r)
	if err != nil {
		return err
	}

	eventLogs := r.EventLogs
	err = model.EventLogs(eventLogs).UpdateToUsed(s.db)
	if err != nil {
		return err
	}

	return nil
}

func (s realTimeService) UpdateResult(req request.UpdateResult) *pkg.Error {
	r := model.Record{ID: req.RecordID}
	err := s.recordExist(&r)
	if err != nil {
		return err
	}

	fieldResultUpdateList := make([]model.FieldResult, 0, len(req.Fields))
	for i := range req.Fields {
		fieldResultUpdateList = append(fieldResultUpdateList, req.Fields[i].Results...)
	}

	result := s.db.Save(&fieldResultUpdateList)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}

	return nil
}

func (s realTimeService) GetResult(r *model.Record) ([]model.Event, int64, *pkg.Error) {
	err := s.recordExist(r)
	if err != nil {
		return nil, 0, err
	}

	e := model.Event{}
	//eventIDs, _ := pkg.Strs{}.Scan(r.Events)
	es := model.Events(r.Events)

	eventIDs, _ := es.ListEventID()
	events, totalRow, err := e.ListWithNewestResult(s.db, eventIDs, r.ID)

	if err != nil {
		return nil, 0, err
	}

	return events, totalRow, nil
}

func (s realTimeService) recordExist(r *model.Record) *pkg.Error {
	err := r.Get(s.db)
	if err != nil {
		if strings.Contains(err.Msg, "record not found") {
			return pkg.NewError(pkg.NotFound, fmt.Sprintf("record with id %s not found", r.ID))
		}
		return err
	}
	return nil
}

func (s realTimeService) Test(wp *worker.Pool, req request.Start) (model.Record, *pkg.Error) {
	// 判断events合不合法
	_, eventNames, err, events := eventsLegitimate(s.db, req.EventIDs)
	if err != nil {
		return model.Record{}, err
	}
	// 判断accounts合不合法
	err = accountsLegitimate(s.db, req.AccountIDs)
	if err != nil {
		return model.Record{}, err
	}

	// 判断project存不存在
	p := model.Project{
		ID: req.Project,
	}
	err = p.Get(s.db)
	if err != nil {
		if strings.Contains(err.Msg, "record not found") {
			return model.Record{}, pkg.NewError(pkg.BadRequest, "project does not exist")
		}
		return model.Record{}, err
	}

	// 创建filter
	filter := jet.Filter{
		Events:  eventNames,
		Project: p.ID,
		UserIDs: req.AccountIDs,
	}
	filterRes, err := filter.Create()
	if err != nil {
		return model.Record{}, err
	}

	if filterRes.Status != "READY" {
		return model.Record{}, pkg.NewError(pkg.ServerError, "调用jet创建filter时出了问题，请稍后再试")
	}

	filter.ID = filterRes.ID
	filter.Status = jet.RECORDING
	err = filter.Update()
	if err != nil {
		return model.Record{}, err
	}

	// 创建record
	u := uuid.New()
	id := strings.ReplaceAll(u.String(), "-", "")

	r := model.Record{
		Name:      "实时埋点测试" + time.Stamp,
		Status:    model.ON,
		Filter:    filter.ID,
		ProjectID: filter.Project,
		ID:        id,
		//Events:    eventsValue,
		Events: events,
	}
	err = r.Create(s.db)
	if err != nil {
		return r, err
	}

	// 2.先尝试用go的标准库time和channel加上全剧的workerPool
	l := os.Getenv("REALTIME_LIMIT")
	limit := pkg.StrTo(l).MustInt()

	go checkLog(limit, r, wp)

	return r, nil
}

func eventsLegitimate(db *gorm.DB, ids []string) ([]string, []string, *pkg.Error, []model.Event) {
	e := model.Event{}
	events, totalRow, err := e.List(db, ids)
	if err != nil {
		return nil, nil, pkg.NewError(pkg.ServerError, err.Error()), nil
	}

	eventNames := make([]string, len(events))
	for _, event := range events {
		eventNames = append(eventNames, event.Name)
	}

	if hasDuplicates(eventNames) {
		return nil, nil, pkg.NewError(pkg.BadRequest, "不能传入名字是一样的event"), nil
	}

	eventIDs := make([]string, len(events))
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

	if totalRow < int64(len(ids)) {
		return nil, nil, pkg.NewError(pkg.BadRequest, "传入的events id有一部分不存在"), nil
	}

	return eventIDs, eventNames, nil, events
}

func accountsLegitimate(db *gorm.DB, ids []string) *pkg.Error {
	a := model.Account{}
	_, totalRow, err := a.GetSome(db, ids)

	if totalRow < int64(len(ids)) {
		return pkg.NewError(pkg.BadRequest, "传入的accounts id有一部分不存在")
	}

	if err != nil {
		return err
	}
	return nil
}

// 每2秒会往task channel里塞一个任务
func checkLog(limit int, r model.Record, wp *worker.Pool) {
	ticker := time.NewTicker(2 * time.Second)
	i := 1
	for range ticker.C {
		initializer.DB.First(&r)
		if i > limit || r.Status == model.OFF {
			r.Status = model.OFF
			initializer.DB.Save(r)

			filter := jet.Filter{
				ID:     r.Filter,
				Status: jet.STOPPED,
			}
			_ = filter.Update()
			ticker.Stop()
			break
		}

		job := task.RealTimeTrackTask{
			Type:   task.Fetch,
			Record: r,
			WP:     wp,
		}

		select {
		case wp.Jobs <- &job:
			fmt.Println("Fetch job has been added to the queue")
		default:
			fmt.Println("Too many jobs")
		}
		i += 1
	}
}

func hasDuplicates(strs []string) bool {
	set := make(map[string]bool)
	for _, str := range strs {
		if set[str] {
			return true
		}
		set[str] = true
	}
	return false
}
