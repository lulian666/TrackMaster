package router

import (
	"TrackMaster/handler"
	"TrackMaster/initializer"
	"TrackMaster/service"
	"github.com/gin-gonic/gin"
)

func AddRealTimeRoutes(g *gin.RouterGroup) {
	realTimeService := service.NewRealTimeService(initializer.DB)
	realTimeHandler := handler.NewRealTimeHandler(realTimeService)

	g.POST("/realTime/start", realTimeHandler.Start)
	g.POST("/realTime/stop", realTimeHandler.Stop)
	g.POST("/realTime/update", realTimeHandler.Update)
	g.GET("/realTime/getLog", realTimeHandler.GetLog)
	g.POST("/realTime/clearLog", realTimeHandler.ClearLog)
	g.POST("/realTime/resetResult", realTimeHandler.UpdateResult)
	g.GET("/realTime/getResult", realTimeHandler.GetResult)

	//测试
	g.POST("/realTime/test", realTimeHandler.Test)
}
