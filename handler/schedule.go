package handler

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/pkg/worker"
	"TrackMaster/service"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	service service.ScheduleService
}

func NewScheduleHandler(s service.ScheduleService) ScheduleHandler {
	return ScheduleHandler{
		service: s,
	}
}

// On
// @Tags schedule
// @Summery set schedule for project
// @Produce json
// @Param project query string true "project id"
// @Success 200 {object} object "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/schedules/on [post]
func (h ScheduleHandler) On(wp *worker.Pool) func(c *gin.Context) {
	return func(c *gin.Context) {
		res := pkg.NewResponse(c)

		// 字段是project，值是id
		projectID := c.Query("project")
		if projectID == "" {
			res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "project cannot be null"))
			return
		}

		p := model.Project{
			ID: projectID,
		}

		schedule, err := h.service.On(&p, wp)
		if err != nil {
			res.ToErrorResponse(err)
		}

		res.ToResponse(schedule)
	}
}

// Off
// @Tags schedule
// @Summery set schedule for project
// @Produce json
// @Param project query string true "project id"
// @Success 200 {object} object "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/schedules/off [post]
func (h ScheduleHandler) Off(c *gin.Context) {
	res := pkg.NewResponse(c)

	// 字段是project，值是id
	projectID := c.Query("project")
	if projectID == "" {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "project cannot be null"))
		return
	}

	p := model.Project{
		ID: projectID,
	}

	schedule, err := h.service.Off(&p)
	if err != nil {
		res.ToErrorResponse(err)
	}

	res.ToResponse(schedule)
}
