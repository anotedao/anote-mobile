package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	db, err := gorm.Open(postgres.Open(conf.DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		log.Println(err)
	}

	if err := db.AutoMigrate(&Miner{}, &IpAddress{}, &KeyValue{}); err != nil {
		panic(err.Error())
	}

	return db
}
