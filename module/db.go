package module

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DBConnectionSSO() (*gorm.DB, error) {
	USER := "app_sso"
	PASS := "Admin123##"
	HOST := "192.168.11.153"
	PORT := "6033"
	DBNAME := "core_sso"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)

	return gorm.Open(mysql.Open(url), &gorm.Config{Logger: newLogger})

}
