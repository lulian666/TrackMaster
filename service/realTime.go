package service

import (
	"TrackMaster/model"
	"TrackMaster/model/request"
	"TrackMaster/pkg"
	"TrackMaster/third_party/jet"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type RealTimeService interface {
	Start(req request.Request) (model.Record, *pkg.Error)
	Stop()
	Update()
	GetLog()
	ClearLog()
	ResetResult()
}

type realTimeService struct {
	db *gorm.DB
}

func NewRealTimeService(db *gorm.DB) RealTimeService {
	return &realTimeService{
		db: db,
	}
}

func (s realTimeService) Start(req request.Request) (model.Record, *pkg.Error) {
	// 判断events存不存在
	e := model.Event{}
	events, totalRow, err := e.ListEventName(s.db, req.EventIDs)

	if totalRow < int64(len(req.EventIDs)) {
		return model.Record{}, pkg.NewError(pkg.BadRequest, "传入的events id有一部分不存在")
	}

	if err != nil {
		return model.Record{}, pkg.NewError(pkg.ServerError, err.Error())
	}

	// 判断accounts存不存在
	a := model.Account{}
	_, totalRow, err = a.GetSome(s.db, req.AccountIDs)

	if totalRow < int64(len(req.AccountIDs)) {
		return model.Record{}, pkg.NewError(pkg.BadRequest, "传入的accounts id有一部分不存在")
	}

	if err != nil {
		return model.Record{}, pkg.NewError(pkg.ServerError, err.Error())
	}

	// 判断project存不存在
	p := model.Project{
		ID: req.Project,
	}
	err = p.Get(s.db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Record{}, pkg.NewError(pkg.BadRequest, "project does not exist")
		}
		return model.Record{}, pkg.NewError(pkg.ServerError, err.Error())
	}

	// 创建filter
	filter := jet.Filter{
		Events:  events,
		Project: p.ID,
		UserIDs: req.AccountIDs,
	}
	filterRes, err := filter.Create()
	if err != nil {
		return model.Record{}, pkg.NewError(pkg.ServerError, err.Error())
	}

	if filterRes.Status != "READY" {
		return model.Record{}, pkg.NewError(pkg.ServerError, "调用jet创建filter时出了问题，请稍后再试")
	}

	filter.ID = filterRes.ID
	filter.Status = jet.RECORDING
	err = filter.Update()
	if err != nil {
		return model.Record{}, pkg.NewError(pkg.ServerError, err.Error())
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
	}
	err = r.Create(s.db)
	if err != nil {
		return r, pkg.NewError(pkg.ServerError, err.Error())
	}

	// todo 创建定时任务
	// 1.先尝试用go的标准库time和channel来实现试试

	return r, nil
}

func (s realTimeService) Stop() {

}

func (s realTimeService) Update() {

}

func (s realTimeService) GetLog() {

}

func (s realTimeService) ClearLog() {

}

func (s realTimeService) ResetResult() {

}
