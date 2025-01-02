package config

import (
	"ApiSup/pkg/mapear/constants"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_SENHA"),
		os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)), &gorm.Config{})
}

var DB *gorm.DB

func Initialize() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal(constants.ERROR_CARREGAMENTO_ENV)
	}

	dbConn, err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	DB = dbConn

	return DB
}

func Close() error {
	sqlDB, err := DB.DB()

	if err != nil {
		return err
	}

	return sqlDB.Close()
}
