package track

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/third_party/jet"
	"TrackMaster/third_party/podcast"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func FetchNewLog(r model.Record) (int, *pkg.Error) {
	fmt.Println("fetch NewLog...")
	logs, err := jet.GetLogs(r.Filter)
	if err != nil {
		return 0, err.WithDetails("something went wrong when fetching logs from jet")
	}

	if len(logs) > 0 {
		fmt.Println("done fetching, clear log...")
		err = jet.ClearLogs(r.Filter)
		if err != nil {
			return 0, err.WithDetails("something went wrong when clearing logs from jet")
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

			// 被测的events都记录在record的Events字段里
			es := model.Events(r.Events)
			eventIDs, _ := es.ListEventID()
			e := model.Event{}
			events, _, err := e.List(initializer.DB, eventIDs)
			if err != nil {
				return 0, err.WithDetails("something went wrong when listing all events from db")
			}

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
				return 0, pkg.NewError(pkg.ServerError, result.Error.Error())
			}
		}

		// 理论上所有收集到的log都是按照events去过滤的，每一条log都是需要存的
		if len(fieldLogCreateList) > 0 {
			result := initializer.DB.Create(fieldLogCreateList)
			if result.Error != nil {
				return 0, pkg.NewError(pkg.ServerError, result.Error.Error())
			}
		}
	}
	return len(logs), nil
}

func TestNewLog(r model.Record) *pkg.Error {
	fmt.Println("testing NewLog...")
	// 这里传入的record应该并没有preload EventLogs
	initializer.DB.Model(r).Preload("EventLogs.FieldLogs").First(&r)
	var iosLogs []model.EventLog
	var androidLogs []model.EventLog
	var otherLogs []model.EventLog

	for _, log := range r.EventLogs {
		if log.Platform == "iOS" && !log.Tested {
			iosLogs = append(iosLogs, log)
		} else if log.Platform == "Android" && !log.Tested {
			androidLogs = append(androidLogs, log)
		} else if !log.Tested {
			otherLogs = append(otherLogs, log)
		}
	}

	if len(iosLogs) > 0 {
		err := testEvent(iosLogs, r)
		if err != nil {
			return err
		}
	}

	if len(androidLogs) > 0 {
		err := testEvent(androidLogs, r)
		if err != nil {
			return err
		}
	}

	if len(otherLogs) > 0 {
		err := testEvent(otherLogs, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func testEvent(eventLogs []model.EventLog, r model.Record) *pkg.Error {
	// 被测的events都记录在record的Events字段里
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
				Result: model.Result{
					IOS:     model.UNTESTED,
					Android: model.UNTESTED,
					Other:   model.UNTESTED,
				},
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

	result := initializer.DB.Create(&eventResultCreateList)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error()).WithDetails("something went wrong creating event results")
	}

	result = initializer.DB.Create(&fieldResultCreateList)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error()).WithDetails("something went wrong creating field results")
	}

	result = initializer.DB.Model(&eventLogUpdateList).Update("tested", true)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error()).WithDetails("something went wrong updating event logs")
	}

	return nil
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
		Where("platform = ? AND field_id = ? AND `key` = ?", fieldLog.Platform, fieldLog.FieldID, fieldLog.Key).
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
