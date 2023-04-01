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
	//err := initializer.DB.AutoMigrate(&model.Project{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}
	//
	//err = initializer.DB.AutoMigrate(&model.Account{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}
	//
	//err = initializer.DB.AutoMigrate(&model.Type{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}
	//
	//err = initializer.DB.AutoMigrate(&model.EnumValue{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}
	//
	//err = initializer.DB.AutoMigrate(&model.Story{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}
	//
	//err = initializer.DB.AutoMigrate(&model.Event{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}
	//
	//err = initializer.DB.AutoMigrate(&model.Type{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}
	//
	//err = initializer.DB.AutoMigrate(&model.EnumValue{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}
	//
	//err = initializer.DB.AutoMigrate(&model.Story{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}
	//
	//err = initializer.DB.AutoMigrate(&model.Event{})
	//if err != nil {
	//	log.Fatal("Error migrating database")
	//}

	err := initializer.DB.AutoMigrate(&model.Field{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	// 上面这个创建表时，联合主键没生效，所以用sql创建了表
	// 后面有空研究下为什么
	//CREATE TABLE `fields` (
	//	`event_id` varchar(191),
	//	`id` varchar(191),
	//	`type` longtext,
	//	`type_id` longtext,
	//	`key` longtext,
	//	`value` longtext,
	//	`description` longtext,
	//	`created_at` datetime(3) NULL,
	//	`updated_at` datetime(3) NULL,
	//	PRIMARY KEY (`id`, `event_id`)
	//);

	log.Println("Migration successful")
}
