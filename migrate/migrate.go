package main

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"log"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectDB()
}

func main() {
	err := initializer.DB.AutoMigrate(&model.Project{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	err = initializer.DB.AutoMigrate(&model.Account{})
	if err != nil {
		log.Fatal("Error migrating database")
	}
	log.Println("Migration successful")
}
