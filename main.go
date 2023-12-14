package main

import (
	"log"

	// _ "net/http/pprof"

	"gopkg.in/macaron.v1"
	"gorm.io/gorm"
)

var conf *Config

var mac *macaron.Macaron

var mon *Monitor

var db *gorm.DB

var pc *PriceClient

var cch *Cache

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	conf = initConfig()

	db = initDb()

	mac = initMacaron()

	mon = initMonitor()

	pc = initPriceClient()

	cch = initCache()

	mac.Run("127.0.0.1", Port)
}
