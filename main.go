package main

import (
	_ "TrackMaster/docs"
	"TrackMaster/initializer"
	"TrackMaster/router"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectDB()
}

// main
// @title TrackMaster
// @version 1.0
func main() {
	fmt.Println("hello")
	r := router.NewRouter()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run()
	if err != nil {
		log.Fatal("Failed to start server, due to: ", err.Error())
	}
}
