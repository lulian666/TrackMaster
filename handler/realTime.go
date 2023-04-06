package handler

import (
	"TrackMaster/model"
	"TrackMaster/model/request"
	"TrackMaster/pkg"
	"TrackMaster/pkg/worker"
	"TrackMaster/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type RealTimeHandler struct {
	service service.RealTimeService
}

func NewRealTimeHandler(s service.RealTimeService) RealTimeHandler {
	return RealTimeHandler{
		service: s,
	}
}

// Start
// 会根据被测的events和users来生成一个record并返回id
// @Tags realTime
// @Summery start record
// @Produce json
// @Param project body string true "project ID"
// @Param accounts body []string true "account IDs"
// @Param events body []string true "event IDs"
// @Success 200 {object} model.Record "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/realTime/start [post]
func (h RealTimeHandler) Start(wp *worker.Pool) func(c *gin.Context) {
	return func(c *gin.Context) {
		res := pkg.NewResponse(c)

		req := request.Start{}
		err := c.ShouldBind(&req)
		if err != nil {
			res.ToErrorResponse(pkg.NewError(pkg.BadRequest, err.Error()))
			return
		}

		if len(req.EventIDs) == 0 {
			res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "至少选择一个event"))
			return
		}

		if len(req.AccountIDs) == 0 {
			res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "至少选择一个account"))
			return
		}

		record, err1 := h.service.Start(wp, req)
		if err1 != nil {
			if strings.Contains(err1.Msg, "record not found") {
				res.ToErrorResponse(pkg.NewError(pkg.NotFound, err.Error()))
				return
			}
			res.ToErrorResponse(err1)
			return
		}

		res.ToResponse(record)
	}

}

// Stop
// 停止收集和测试
// @Tags realTime
// @Summery stop record
// @Produce json
// @Param record query string true "record ID"
// @Success 200 {object} model.Record "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/realTime/stop [post]
func (h RealTimeHandler) Stop(c *gin.Context) {
	res := pkg.NewResponse(c)
	recordID := c.Query("record")
	if recordID == "" {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "record required in query"))
		return
	}

	r := model.Record{
		ID: recordID,
	}

	err := h.service.Stop(&r)
	if err != nil {
		res.ToErrorResponse(err)
		return
	}

	res.ToResponse(r)
}

// Update
// @Tags realTime
// @Summery update record
// @Produce json
// @Param record query string true "record ID"
// @Param accounts body []string false "account IDs"
// @Param events body []string false "event IDs"
// @Success 200 {object} model.Record "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/realTime/update [post]
func (h RealTimeHandler) Update(c *gin.Context) {
	res := pkg.NewResponse(c)
	recordID := c.Query("record")
	if recordID == "" {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "record required in query"))
		return
	}

	req := request.Update{}
	err := c.ShouldBind(&req)
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, err.Error()))
		return
	}

	r := model.Record{
		ID: recordID,
	}
	err1 := h.service.Update(&r, req)
	if err1 != nil {
		res.ToErrorResponse(err1)
		return
	}
	res.ToResponse(r)
}

// GetLog
// 读取record下的log，按创建时间倒叙排
// @Tags realTime
// @Summery get log in record
// @Produce json
// @Param record query string true "record ID"
// @Success 200 {object} model.SwaggerEventLogs "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/realTime/getLog [get]
func (h RealTimeHandler) GetLog(c *gin.Context) {
	res := pkg.NewResponse(c)
	recordID := c.Query("record")
	if recordID == "" {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "record required in query"))
		return
	}

	r := model.Record{
		ID: recordID,
	}

	eventLogs, totalRow, err := h.service.GetLog(&r)
	if err != nil {
		res.ToErrorResponse(err)
		return
	}

	res.ToResponseList(eventLogs, totalRow)
}

// ClearLog
// 清除已经收集的log
// @Tags realTime
// @Summery clear log in record
// @Produce json
// @Param record query string true "record ID"
// @Success 200 {object} object "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/realTime/clearLog [post]
func (h RealTimeHandler) ClearLog(c *gin.Context) {
	res := pkg.NewResponse(c)
	recordID := c.Query("record")
	if recordID == "" {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "record required in query"))
		return
	}

	r := model.Record{
		ID: recordID,
	}

	err := h.service.ClearLog(&r)
	if err != nil {
		res.ToErrorResponse(err)
		return
	}

	res.ToResponse(nil)
}

// UpdateResult
// @Tags realTime
// @Summery reset test result in record
// @Produce json
// @Param record body string true "record ID"
// @Param fields body model.Fields true "field IDs"
// @Success 200 {object} object "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/realTime/updateResult [post]
func (h RealTimeHandler) UpdateResult(c *gin.Context) {
	res := pkg.NewResponse(c)
	req := request.UpdateResult{}
	err := c.ShouldBind(&req)
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, err.Error()))
		return
	}

	err1 := h.service.UpdateResult(req)
	if err1 != nil {
		res.ToErrorResponse(err1)
		return
	}
	res.ToResponse(nil)
}

// GetResult
// @Tags realTime
// @Summery get test result in record
// @Produce json
// @Param record query string true "record ID"
// @Success 200 {object} model.SwaggerEvents "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/realTime/getResult [get]
func (h RealTimeHandler) GetResult(c *gin.Context) {
	res := pkg.NewResponse(c)
	recordID := c.Query("record")
	if recordID == "" {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "record required in query"))
		return
	}

	r := model.Record{
		ID: recordID,
	}
	events, totalRow, err := h.service.GetResult(&r)
	if err != nil {
		res.ToErrorResponse(err)
		return
	}

	res.ToResponseList(events, totalRow)
}

func (h RealTimeHandler) Test(wp *worker.Pool) func(c *gin.Context) {
	return func(c *gin.Context) {
		res := pkg.NewResponse(c)
		req := request.Start{}
		err := c.ShouldBind(&req)
		if err != nil {
			res.ToErrorResponse(pkg.NewError(pkg.BadRequest, err.Error()))
			return
		}

		if len(req.EventIDs) == 0 {
			res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "至少选择一个event"))
			return
		}

		if len(req.AccountIDs) == 0 {
			res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "至少选择一个account"))
			return
		}

		record, err := h.service.Start(wp, req)
		if err != nil {
			fmt.Println("error")
		}

		res.ToResponse(record)
	}
}
