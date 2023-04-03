package handler

import (
	"TrackMaster/model"
	"TrackMaster/model/request"
	"TrackMaster/pkg"
	"TrackMaster/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
// @Router /api/v1/realTime/start [post]
func (h RealTimeHandler) Start(c *gin.Context) {
	res := pkg.NewResponse(c)

	req := request.Request{}
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

	record, err1 := h.service.Start(req)
	if err1 != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.ToErrorResponse(pkg.NewError(pkg.NotFound, err.Error()))
			return
		}
		res.ToErrorResponse(err1)
		return
	}

	res.ToResponse(record)
}

func (h RealTimeHandler) Stop(c *gin.Context) {

}

func (h RealTimeHandler) Update(c *gin.Context) {

}

func (h RealTimeHandler) GetLog(c *gin.Context) {

}

func (h RealTimeHandler) ClearLog(c *gin.Context) {

}

func (h RealTimeHandler) ResetResult(c *gin.Context) {

}

func (h RealTimeHandler) Test(c *gin.Context) {
	res := pkg.NewResponse(c)

	r := model.Record{
		ID: "dd6fd66c457e435894322b1fc7ca3496",
	}

	h.service.Test(r)
	res.ToResponse(nil)
}
