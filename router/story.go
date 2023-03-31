package router

import (
	"TrackMaster/handler"
	"TrackMaster/initializer"
	"TrackMaster/service"
	"github.com/gin-gonic/gin"
)

func AddStoryRoutes(g *gin.RouterGroup) {
	storyService := service.NewStoryService(initializer.DB)
	storyHandler := handler.NewStoryHandler(storyService)

	g.POST("/stories/sync", storyHandler.Sync)
}
