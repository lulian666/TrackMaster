package main

import (
	_ "TrackMaster/docs"
	"TrackMaster/initializer"
	"TrackMaster/pkg"
	"TrackMaster/pkg/worker"
	"TrackMaster/router"
	"TrackMaster/third_party/slack"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"strings"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectDB()
}

// main
// @title TrackMaster
// @version 2.0
func main() {
	errorCh := make(chan *pkg.Error, worker.MaxQueue)
	wp := worker.NewWorkerPool(errorCh)
	wp.Start()

	r := router.NewRouter(wp)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go func() {
		err := r.Run()
		if err != nil {
			log.Fatal("Failed to start server, due to: ", err.Error())
		}
	}()

	// 监测错误
	for err := range errorCh {
		fmt.Printf("Error occurred: %v\n", err.Error())
		fmt.Printf("Error details: %v\n", err.Details)
		details := strings.Join(err.Details, ";")
		_ = slack.SendMessage(err.Error() + "\n" + details)
	}

}
