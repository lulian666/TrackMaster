package service

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"TrackMaster/model/request"
	"TrackMaster/pkg"
	"TrackMaster/third_party/jet"
	"TrackMaster/third_party/podcast"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type RealTimeService interface {
	Start(req request.Start) (model.Record, *pkg.Error)
	Stop(r *model.Record) *pkg.Error
	Update(r *model.Record, req request.Update) *pkg.Error
	GetLog(r *model.Record) ([]model.EventLog, int64, *pkg.Error)
	ClearLog(r *model.Record) *pkg.Error
	ResetResult()
	GetResult(r *model.Record) ([]model.Event, int64, *pkg.Error)
	Test(r model.Record)
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

func (s realTimeService) Start(req request.Start) (model.Record, *pkg.Error) {
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
	err1 := p.Get(s.db)
	if err1 != nil {
		if errors.Is(err1, gorm.ErrRecordNotFound) {
			return model.Record{}, pkg.NewError(pkg.BadRequest, "project does not exist")
		}
		return model.Record{}, pkg.NewError(pkg.ServerError, err1.Error())
	}

	// 创建filter
	filter := jet.Filter{
		Events:  eventNames,
		Project: p.ID,
		UserIDs: req.AccountIDs,
	}
	filterRes, err1 := filter.Create()
	if err1 != nil {
		return model.Record{}, pkg.NewError(pkg.ServerError, err1.Error())
	}

	if filterRes.Status != "READY" {
		return model.Record{}, pkg.NewError(pkg.ServerError, "调用jet创建filter时出了问题，请稍后再试")
	}

	filter.ID = filterRes.ID
	filter.Status = jet.RECORDING
	err1 = filter.Update()
	if err1 != nil {
		return model.Record{}, pkg.NewError(pkg.ServerError, err1.Error())
	}

	// 创建record
	u := uuid.New()
	id := strings.ReplaceAll(u.String(), "-", "")
	//eventsValue, _ := pkg.Strs(eventIDs).Value()

	r := model.Record{
		Name:      "实时埋点测试" + time.Stamp,
		Status:    model.ON,
		Filter:    filter.ID,
		ProjectID: filter.Project,
		ID:        id,
		//Events:    eventsValue,
		Events: events,
	}
	err1 = r.Create(s.db)
	if err1 != nil {
		return r, pkg.NewError(pkg.ServerError, err1.Error())
	}

	// todo 创建存log和测log的任务
	// 1.先尝试用go的标准库time和channel来实现试试
	logCh := make(chan int)
	// todo 需要做错误处理
	go checkLog(2000, logCh, r) // todo 后面写成可配置
	go testLog(logCh, r)

	return r, nil
}

func (s realTimeService) Stop(r *model.Record) *pkg.Error {
	err := r.Update(s.db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NewError(pkg.NotFound, fmt.Sprintf("record with id %s not found", r.ID))
		}
		return pkg.NewError(pkg.ServerError, err.Error())
	}

	filter := jet.Filter{
		ID:     r.Filter,
		Status: jet.STOPPED,
	}
	_ = filter.Update()

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

	err1 := filter.Update()
	if err1 != nil {
		return pkg.NewError(pkg.ServerError, err1.Error())
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, pkg.NewError(pkg.NotFound, fmt.Sprintf("record with id %s not found", r.ID))
		}
		return nil, 0, pkg.NewError(pkg.ServerError, err.Error())
	}

	return logs, totalRow, nil
}

func (s realTimeService) ClearLog(r *model.Record) *pkg.Error {
	err := s.recordExist(r)
	if err != nil {
		return err
	}

	eventLogs := r.EventLogs
	err1 := model.EventLogs(eventLogs).UpdateToUsed(s.db)
	if err1 != nil {
		return pkg.NewError(pkg.ServerError, err1.Error())
	}

	return nil
}

func (s realTimeService) ResetResult() {

}

func (s realTimeService) GetResult(r *model.Record) ([]model.Event, int64, *pkg.Error) {
	err := s.recordExist(r)
	if err != nil {
		return nil, 0, err
	}

	e := model.Event{}
	//eventIDs, _ := pkg.Strs{}.Scan(r.Events)
	es := model.Events(r.Events)

	fmt.Println("events 数量：", len(r.Events))
	eventIDs, _ := es.ListEventID()
	events, totalRow, err1 := e.ListWithNewestResult(s.db, eventIDs, r.ID)

	if err1 != nil {
		return nil, 0, pkg.NewError(pkg.ServerError, err.Error())
	}

	// 批量查询
	//eventResults := model.EventResults(make([]model.EventResult, 0, totalRow))
	//err = eventResults.Get(s.db, *r, events)
	//if err != nil {
	//	return nil, 0, err
	//}

	// 一一对应

	return events, totalRow, nil
}

func (s realTimeService) recordExist(r *model.Record) *pkg.Error {
	err := r.Get(s.db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NewError(pkg.NotFound, fmt.Sprintf("record with id %s not found", r.ID))
		}
		return pkg.NewError(pkg.ServerError, err.Error())
	}
	return nil
}

func (s realTimeService) Test(r model.Record) {
	s.db.First(&r)
	count := fetchNewLog(r)
	fmt.Println("count:", count)
}

func eventsLegitimate(db *gorm.DB, ids []string) ([]string, []string, *pkg.Error, []model.Event) {
	e := model.Event{}
	events, totalRow, err := e.List(db, ids)
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

	if err != nil {
		return nil, nil, pkg.NewError(pkg.ServerError, err.Error()), nil
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
		return pkg.NewError(pkg.ServerError, err.Error())
	}
	return nil
}

// 每2秒会往channel里写一个数（如果有）
// 直到写了limit次以后停止，并且关闭channel
func checkLog(limit int, logCh chan<- int, r model.Record) {
	ticker := time.NewTicker(2 * time.Second)
	i := 1
	for range ticker.C {
		initializer.DB.First(&r)
		if i > limit || r.Status == model.OFF {
			ticker.Stop()
			r.Status = model.OFF
			initializer.DB.Save(r)
			break
		}

		// todo 如果这个函数出错了，别block在这
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

// 真正去取存数据的函数
func fetchNewLog(r model.Record) int {
	fmt.Println("fetch NewLog...")
	logs, _ := jet.GetLogs(r.Filter)
	if len(logs) > 0 {
		fmt.Println("done fetching, clear log...")
		_ = jet.ClearLogs(r.Filter)

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

			// 被测的events都记录在record的Events字段里
			//eventIDs, _ := pkg.Strs{}.Scan(r.Events)
			es := model.Events(r.Events)
			eventIDs, _ := es.ListEventID()
			e := model.Event{}
			events, _, _ := e.List(initializer.DB, eventIDs)

			for j := range events {
				if el.Name == events[j].Name {
					el.EventID = events[j].ID
					completeContent := make(map[string]interface{})

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
						v, ok := log.Get(field.Key, logs[i].LogStr)
						if ok && v != "" {
							fieldLog.Value = v

							// 拿content
							if strings.HasSuffix(field.Key, "id") {
								contentID := v
								keys := strings.Split(field.Key, ".")
								contentTypeKey := keys[0] + ".type"
								contentType, _ := log.Get(contentTypeKey, logs[i].LogStr)
								content, _ := podcast.GetContentByTypeAndID(contentType, contentID)
								completeContent[keys[0]] = content
							}
						}

						fieldLogCreateList = append(fieldLogCreateList, fieldLog)
					}

					content, _ := json.Marshal(completeContent)
					el.Content = string(content)
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
	}
	return len(logs)
}

// 从channel中读数，直到channel被关闭
func testLog(logCh <-chan int, r model.Record) {
	for count := range logCh {
		// this loop closes when channel is closed
		fmt.Printf("reading %d 条新log...\n", count)
		testNewLog(count, r)
	}
	fmt.Println("channel closed")
}

// 真正去测试数据的函数
// count 新增的event log数量
// r 此次被测的record
func testNewLog(count int, r model.Record) {
	fmt.Println("testing NewLog...")
	// 这里传入的record应该并没有preload EventLogs
	initializer.DB.Model(r).Preload("EventLogs.FieldLogs").First(&r)
	iosLogs := make([]model.EventLog, 0, count)
	androidLogs := make([]model.EventLog, 0, count)
	otherLogs := make([]model.EventLog, 0, count)

	for _, log := range r.EventLogs {
		if log.Platform == "iOS" && !log.Tested {
			iosLogs = append(iosLogs, log)
		} else if log.Platform == "Android" && !log.Tested {
			androidLogs = append(androidLogs, log)
		} else if !log.Tested {
			otherLogs = append(otherLogs, log)
		}
	}

	// 三个平台的log一起测
	if len(iosLogs) > 0 {
		go testEvent(iosLogs, r)
	}

	if len(androidLogs) > 0 {
		go testEvent(androidLogs, r)
	}

	if len(otherLogs) > 0 {
		go testEvent(otherLogs, r)
	}

}

func testEvent(eventLogs []model.EventLog, r model.Record) {
	// 被测的events都记录在record的Events字段里
	//eventIDs, _ := pkg.Strs{}.Scan(r.Events)
	es := model.Events(r.Events)
	eventIDs, _ := es.ListEventID()
	e := model.Event{}
	events, _, _ := e.List(initializer.DB, eventIDs)

	eventResultCreateList := make([]model.EventResult, 0, len(eventLogs))
	fieldResultCreateList := make([]model.FieldResult, 0, len(eventLogs))
	eventLogUpdateList := make([]model.EventLog, 0, len(eventLogs))
	fieldLogUpdateList := make([]model.FieldLog, 0, len(eventLogs)*6)
	for _, eventLog := range eventLogs {
		// 创建event测试结构
		u1 := uuid.New()
		id1 := strings.ReplaceAll(u1.String(), "-", "")
		eventResult := model.EventResult{
			RecordID: r.ID,
			EventID:  eventLog.EventID,
			ID:       id1,
		}

		eventResultCreateList = append(eventResultCreateList, eventResult)
		event, _ := model.Events(events).FindByID(eventLog.EventID)

		for _, fieldLog := range eventLog.FieldLogs {
			// 创建field测试结果
			u2 := uuid.New()
			id2 := strings.ReplaceAll(u2.String(), "-", "")
			fieldResult := model.FieldResult{
				RecordID: r.ID,
				FieldID:  fieldLog.FieldID,
				ID:       id2,
			}

			field, _ := model.Fields(event.Fields).FindByID(fieldLog.FieldID)

			// 每一条field log挨个测试
			// 先要找到这个field log对应的field的value
			value := field.Value
			if strings.Contains(value, "|") {
				if enumFieldCheck(field, "|", fieldLog) {
					setFieldResult(fieldLog, &fieldResult, model.SUCCESS)
				} else {
					setFieldResult(fieldLog, &fieldResult, model.FAIL)
				}
			} else if strings.Contains(value, ",") {
				if enumFieldCheck(field, ",", fieldLog) {
					setFieldResult(fieldLog, &fieldResult, model.SUCCESS)
				} else {
					setFieldResult(fieldLog, &fieldResult, model.FAIL)
				}
			} else if strings.Contains(value, "/") {
				if enumFieldCheck(field, "/", fieldLog) {
					setFieldResult(fieldLog, &fieldResult, model.SUCCESS)
				} else {
					setFieldResult(fieldLog, &fieldResult, model.FAIL)
				}
			} else {
				if fieldLog.Value == value {
					setFieldResult(fieldLog, &fieldResult, model.SUCCESS)
				} else if fieldLog.Value != "not found" &&
					(strings.Contains(field.Key, "id") ||
						strings.Contains(field.Key, "content") ||
						field.Value == "") {
					setFieldResult(fieldLog, &fieldResult, model.UNCERTAIN)
				} else {
					setFieldResult(fieldLog, &fieldResult, model.FAIL)
				}
			}
			// 标记测过的event log和field log
			fieldLogUpdateList = append(fieldLogUpdateList, fieldLog)
			fieldResultCreateList = append(fieldResultCreateList, fieldResult)
		}

		eventLogUpdateList = append(eventLogUpdateList, eventLog)
	}
	initializer.DB.Create(&eventResultCreateList)
	initializer.DB.Create(&fieldResultCreateList)
	initializer.DB.Model(&eventLogUpdateList).Update("tested", true)
}

func setEventResult() {

}

// 此处只设置，不修改
func setFieldResult(fieldLog model.FieldLog, fieldResult *model.FieldResult, res model.TestResult) {
	// 先区分是哪端的打点
	switch fieldLog.Platform {
	case "iOS":
		fieldResult.IOS = res
	case "Android":
		fieldResult.Android = res
	default:
		fieldResult.Other = res
	}

}

func enumFieldCheck(field model.Field, sep string, fieldLog model.FieldLog) bool {
	values := strings.Split(field.Value, sep)
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}
	var allFieldLogs []string
	initializer.DB.Model(model.FieldLog{}).
		Where("platform = ?", fieldLog.Platform).
		Where("field_id = ?", fieldLog.FieldID).
		Where("key = ?", fieldLog.Key).
		Pluck("value", &allFieldLogs)
	return stringArrayEqual(values, allFieldLogs)
}

func stringArrayEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, s := range a {
		found := false
		for _, t := range b {
			if s == t {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
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
