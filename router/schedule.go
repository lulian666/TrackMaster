package router

import (
	"TrackMaster/handler"
	"TrackMaster/initializer"
	"TrackMaster/pkg/worker"
	"TrackMaster/service"
	"github.com/gin-gonic/gin"
)

func AddScheduleRouter(g *gin.RouterGroup, wp *worker.Pool) {
	scheduleService := service.NewScheduleService(initializer.DB)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)

	g.POST("/schedules/on", scheduleHandler.On(wp))
	g.POST("/schedules/off", scheduleHandler.Off)
}
