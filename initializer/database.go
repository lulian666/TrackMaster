package initializer

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	// connectURL := os.Getenv("DB_URL")
	user := os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	connectURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	DB, err = gorm.Open(mysql.Open(connectURL), &gorm.Config{})

	if err != nil {
		log.Println("Failed to connect to database")
		log.Fatal("connect url is: ", connectURL)
	}
}
