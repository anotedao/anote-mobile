package main

import (
	"log"

	// _ "net/http/pprof"

	"gopkg.in/macaron.v1"
	"gorm.io/gorm"
)

var conf *Config

var m *macaron.Macaron

var mon *Monitor

var db *gorm.DB

var pc *PriceClient

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	conf = initConfig()

	db = initDb()

	m = initMacaron()

	mon = initMonitor()

	pc = initPriceClient()

	// val := "%d%s__239395"
	// dataTransaction("3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p", &val, nil, nil)

	m.Run("127.0.0.1", Port)
}
