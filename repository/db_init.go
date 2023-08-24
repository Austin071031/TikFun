package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
  "os"
)

var db *gorm.DB



func Init() error {
  MYSQL_HOST := os.Getenv("MYSQL_HOST")
  MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
  MYSQL_PORT := os.Getenv("MYSQL_PORT")
  MYSQL_USER := os.Getenv("MYSQL_USER")

  dsn := MYSQL_USER + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_HOST + ":" + MYSQL_PORT + ")/tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}