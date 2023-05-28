package main

import (
	"log"

	"gopkg.in/macaron.v1"
	"gorm.io/gorm"
)

var conf *Config

var m *macaron.Macaron

var mon *Monitor

var db *gorm.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	db = initDb()

	m = initMacaron()

	mon = initMonitor()

	// val := "%d%s__219079"
	// dataTransaction("3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p", &val, nil, nil)

	m.Run("127.0.0.1", Port)
}
