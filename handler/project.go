package handler

import (
	"TrackMaster/pkg"
	"TrackMaster/service"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service service.ProjectService
}

func NewProjectHandler(s service.ProjectService) ProjectHandler {
	return ProjectHandler{
		service: s,
	}
}

// Sync
// @Tags project
// @Summery sync project
// @Produce json
// @Success 200 {object} object "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v1/projects/sync [post]
// 更新本地的project表，让他和jet上的保持同步
func (h ProjectHandler) Sync(c *gin.Context) {
	res := pkg.NewResponse(c)
	err := h.service.SyncProject()
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.ServerError, err.Error()))
		return
	}
	res.ToResponse(nil)
}

// List
// @Tags project
// @Summery list project
// @Produce json
// @Success 200 {array} model.Projects "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v1/projects [get]
func (h ProjectHandler) List(c *gin.Context) {
	res := pkg.NewResponse(c)
	projects, err := h.service.ListProject()
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.ServerError, err.Error()))
		return
	}
	res.ToResponseList(projects, int64(len(projects)))
}
