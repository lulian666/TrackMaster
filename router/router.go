package router

import (
	"TrackMaster/pkg/worker"
	"github.com/gin-gonic/gin"
)

func NewRouter(wp *worker.Pool) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiV1 := r.Group("api/v2")

	AddProjectRoutes(apiV1)
	AddAccountRoutes(apiV1)
	AddStoryRoutes(apiV1)
	AddRealTimeRoutes(apiV1, wp)

	return r
}
