package router

import (
	"TrackMaster/handler"
	"TrackMaster/initializer"
	"TrackMaster/middleware"
	"TrackMaster/model"
	"TrackMaster/service"
	"github.com/gin-gonic/gin"
)

func AddAccountRoutes(g *gin.RouterGroup) {
	accountService := service.NewAccountService(initializer.DB)
	accountHandler := handler.NewAccountHandler(accountService)

	g.POST("/accounts", middleware.Validator(&model.Account{}), accountHandler.Create)
	g.GET("/accounts", accountHandler.List)
	g.DELETE("/accounts/:id", accountHandler.Delete)
}
