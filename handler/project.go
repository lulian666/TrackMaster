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

// SyncProject
// @Tags project
// @Summery sync project
// @Produce json
// @Success 200 {object} object "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v1/sync-projects [post]
// 更新本地的project表，让他和jet上的保持同步
func (h ProjectHandler) SyncProject(c *gin.Context) {
	res := pkg.NewResponse(c)
	err := h.service.SyncProject()
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.ServerError, err.Error()))
	}
	res.ToResponse(nil)
}

// ListProjects
// @Tags project
// @Summery list project
// @Produce json
// @Success 200 {array} model.SwaggerProjects "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v1/projects [get]
func (h ProjectHandler) ListProjects(c *gin.Context) {
	res := pkg.NewResponse(c)
	projects, err := h.service.ListProject()
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.ServerError, err.Error()))
	}
	res.ToResponseList(projects, int64(len(projects)))
}
