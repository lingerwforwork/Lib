package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectMysqlDataBase() (*gorm.DB, error) {
	user := os.Getenv("MYSQL_USER")
	if len(user) <= 0 {
		return nil, fmt.Errorf("MYSQL_USER is not set")
	}
	password := os.Getenv("MYSQL_PASSWORD")
	if len(password) <= 0 {
		return nil, fmt.Errorf("MYSQL_PASSWORD is not set")
	}
	host := os.Getenv("MYSQL_HOST")
	if len(host) <= 0 {
		return nil, fmt.Errorf("MYSQL_HOST is not set")
	}
	port := os.Getenv("MYSQL_PORT")
	if len(port) <= 0 {
		return nil, fmt.Errorf("MYSQL_PORT is not set")
	}
	db := os.Getenv("MYSQL_DATABASE")
	if len(db) <= 0 {
		return nil, fmt.Errorf("MYSQL_DATABASE is not set")
	}
	config := gorm.Config{}
	if os.Getenv("ENVIRONMENT") == "PROD" || os.Getenv("ENVIRONMENT") == "STG" {
		config.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Error,
			},
		)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", user, password, host, port, db)
	gromDB, err := gorm.Open(mysql.Open(dsn), &config)
	if err != nil {
		return nil, err
	}
	return gromDB, nil
}
