package router

import (
	"TrackMaster/pkg/worker"
	"github.com/gin-gonic/gin"
)

func NewRouter(wp *worker.Pool) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiV2 := r.Group("api/v2")

	AddProjectRoutes(apiV2)
	AddAccountRoutes(apiV2)
	AddStoryRoutes(apiV2)
	AddRealTimeRoutes(apiV2, wp)
	AddScheduleRouter(apiV2, wp)

	return r
}
