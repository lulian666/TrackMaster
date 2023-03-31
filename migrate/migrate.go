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

	err = initializer.DB.AutoMigrate(&model.Type{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	err = initializer.DB.AutoMigrate(&model.EnumValue{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	err = initializer.DB.AutoMigrate(&model.Story{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	err = initializer.DB.AutoMigrate(&model.Event{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	err = initializer.DB.AutoMigrate(&model.Type{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	err = initializer.DB.AutoMigrate(&model.EnumValue{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	err = initializer.DB.AutoMigrate(&model.Story{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	err = initializer.DB.AutoMigrate(&model.Event{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	//err := initializer.DB.AutoMigrate(&model.Field{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}

	log.Println("Migration successful")
}
