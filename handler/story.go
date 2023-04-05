package handler

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
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
// @Router /api/v2/stories/sync [post]
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
		if strings.Contains(err.Msg, "record not found") {
			res.ToErrorResponse(pkg.NewError(pkg.NotFound, "project does not exist"))
			return
		}
		res.ToErrorResponse(err)
		return
	}

	res.ToResponse(nil)
}

// List
// @Tags story
// @Summery list story
// @Produce json
// @Param page query string false "page, default 10"
// @Param pageSize query string false "page size, default 10"
// @Success 200 {object} model.Stories "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/stories [get]
func (h StoryHandler) List(c *gin.Context) {
	res := pkg.NewResponse(c)
	pager := pkg.Pager{
		Page:     pkg.GetPage(c),
		PageSize: pkg.GetPageSize(c),
	}

	// 按project取
	projectID := c.Query("project")
	if projectID == "" {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "project cannot be null"))
		return
	}

	s := model.Story{
		ProjectID: projectID,
	}
	stories, totalRow, err := h.service.ListStory(&s, pager)
	if err != nil {
		if !strings.Contains(err.Msg, "record not found") {
			res.ToErrorResponse(err)
			return
		}
		res.ToErrorResponse(pkg.NewError(pkg.NotFound, fmt.Sprintf("project with id %s does not exist", projectID)))
		return
	}
	res.ToResponseList(stories, totalRow)
}

// Get
// @Tags story
// @Summery get story
// @Produce json
// @Param id path string true "story id"
// @Success 200 {object} model.Story "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/stories/{id} [get]
func (h StoryHandler) Get(c *gin.Context) {
	res := pkg.NewResponse(c)
	id := c.Param("id")
	s := model.Story{
		ID: id,
	}

	err := h.service.GetStory(&s)
	if err != nil {
		if !strings.Contains(err.Msg, "record not found") {
			res.ToErrorResponse(err)
			return
		}
		res.ToErrorResponse(pkg.NewError(pkg.NotFound, "story does not exist"))
		return
	}

	res.ToResponse(s)
}
