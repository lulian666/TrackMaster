package router

import (
	"TrackMaster/handler"
	"TrackMaster/initializer"
	"TrackMaster/service"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	projectService := service.NewProjectService(initializer.DB)
	projectHandler := handler.NewProjectHandler(projectService)

	apiV1 := r.Group("api/v1")
	{
		apiV1.POST("/sync-projects", projectHandler.SyncProject)
		apiV1.GET("/projects", projectHandler.ListProjects)
	}

	return r
}
