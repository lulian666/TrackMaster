package handler

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StoryHandler struct {
	service service.StoryService
}

func NewStoryHandler(s service.StoryService) StoryHandler {
	return StoryHandler{
		service: s,
	}
}

// Sync
// @Tags story
// @Summery sync story
// @Produce json
// @Param project query string true "project id"
// @Success 200 {array} model.Story "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v1/stories/sync [post]
func (h StoryHandler) Sync(c *gin.Context) {
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

	err := h.service.SyncStory(&p)
	if err != nil {
		// 有可能是没找到project的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "project does not exist"))
			return
		}
		res.ToErrorResponse(pkg.NewError(pkg.ServerError, err.Error()))
		return
	}

	res.ToResponse(nil)
}
