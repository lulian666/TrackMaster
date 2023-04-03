package service

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"TrackMaster/model/request"
	"TrackMaster/pkg"
	"TrackMaster/third_party/jet"
	"errors"
	"fmt"
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
	Test(r model.Record)
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
	events, totalRow, err := e.List(s.db, req.EventIDs)
	eventNames := make([]string, len(events))
	for _, event := range events {
		eventNames = append(eventNames, event.Name)
	}

	// todo event name 不能重复

	eventIDs := make([]string, len(events))
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

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
		Events:  eventNames,
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
	eventsValue, _ := pkg.Strs(eventIDs).Value()

	r := model.Record{
		Name:      "实时埋点测试" + time.Stamp,
		Status:    model.ON,
		Filter:    filter.ID,
		ProjectID: filter.Project,
		ID:        id,
		Events:    eventsValue,
	}
	err = r.Create(s.db)
	if err != nil {
		return r, pkg.NewError(pkg.ServerError, err.Error())
	}

	// todo 创建存log和测log的任务
	// 1.先尝试用go的标准库time和channel来实现试试
	logCh := make(chan int)
	// todo 需要做错误处理
	go checkLog(2000, logCh, r) // todo 后面写成可配置
	go testLog(logCh)

	return r, nil
}

// 每2秒会往channel里写一个数
// 直到写了limit次以后停止，并且关闭channel
func checkLog(limit int, logCh chan<- int, r model.Record) {
	ticker := time.NewTicker(2 * time.Second)
	i := 1
	for range ticker.C {
		if i > limit {
			ticker.Stop()
			r.Status = model.OFF
			initializer.DB.Save(r)
			break
		}

		count := fetchNewLog(r)

		if count > 0 {
			fmt.Printf("正在写入channel，一共获取了%d条新Log\n", count)
			logCh <- count // 这里应该拿到新log了就往channel里面写
		}

		i += 1
	}
	fmt.Println("channel closed")
	close(logCh)
}

func fetchNewLog(r model.Record) int {
	fmt.Println("fetchNewLog...")
	logs, _ := jet.GetLogs(r.Filter)
	if len(logs) > 0 {
		fmt.Println("done fetching, clear log...")
		_ = jet.ClearLogs(r.Filter)
	}

	// 存log
	// 找到event和eventLog的对应关系
	logCreateList := make([]model.EventLog, 0, len(logs))
	fieldLogCreateList := make([]model.FieldLog, 0, len(logs)*10)
	for i := range logs {
		el := model.EventLog{
			RecordID: r.ID,
			ID:       logs[i].ID,
			Name:     logs[i].Event,
			UserID:   logs[i].UserID,
			Platform: logs[i].Log.OS,
			Raw:      logs[i].LogStr,
		}

		// todo Content 可以放到测试log的时候再去更新

		// 被测的events都记录在record的Events字段里
		eventIDs, _ := pkg.Strs{}.Scan(r.Events)
		e := model.Event{}
		events, _, _ := e.List(initializer.DB, eventIDs)

		for j := range events {
			if el.Name == events[j].Name {
				el.EventID = events[j].ID

				// 遍历events[j]里的fields
				for _, field := range events[j].Fields {
					// 有一个field(需求)，就要创建一个fieldLog(结果)
					// 结果可以为空，但必须有记录
					u := uuid.New()
					id := strings.ReplaceAll(u.String(), "-", "")
					fieldLog := model.FieldLog{
						EventLogID: el.ID,
						FieldID:    field.ID,
						ID:         id,
						Key:        field.Key,
						Value:      "not found", // 默认值，如果找到了就填充进去
						Platform:   el.Platform,
					}

					// 如果是是app的打点，传上来的key是xx.xx格式
					// 如果是前端打点，传上来的key是xx$$xx格式
					log := logs[i].Log
					v, ok := log.Get(field.Key)
					if ok && v != "" {
						fieldLog.Value = v
					}
					fieldLogCreateList = append(fieldLogCreateList, fieldLog)
				}
			}
		}

		logCreateList = append(logCreateList, el)
	}

	if len(logCreateList) > 0 {
		result := initializer.DB.Create(logCreateList)
		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
	}

	// 理论上所有收集到的log都是按照events去过滤的，每一条log都是需要存的
	if len(fieldLogCreateList) > 0 {
		result := initializer.DB.Create(fieldLogCreateList)
		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
	}

	return len(logCreateList)
}

// 从channel中读数，直到channel被关闭
func testLog(logCh <-chan int) {
	for i := range logCh {
		// this loop closes when channel is closed
		fmt.Println("reading...", i)
	}
	fmt.Println("channel closed")
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

func (s realTimeService) Test(r model.Record) {
	s.db.First(&r)
	count := fetchNewLog(r)
	fmt.Println("count:", count)
}
