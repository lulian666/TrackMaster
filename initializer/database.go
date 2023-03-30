package initializer

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	connectURL := os.Getenv("DB_URL")
	DB, err = gorm.Open(mysql.Open(connectURL), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
}
