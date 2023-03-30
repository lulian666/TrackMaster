package router

import (
	"TrackMaster/handler"
	"TrackMaster/initializer"
	"TrackMaster/service"
	"github.com/gin-gonic/gin"
)

func AddProjectRoutes(g *gin.RouterGroup) {
	projectService := service.NewProjectService(initializer.DB)
	projectHandler := handler.NewProjectHandler(projectService)

	g.POST("/projects/sync", projectHandler.Sync)
	g.GET("/projects", projectHandler.List)

}
