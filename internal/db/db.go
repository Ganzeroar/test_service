package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(url string) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{QueryFields: true})

	if err != nil {
		log.Fatalln(err)
	}

	DB = db
}
