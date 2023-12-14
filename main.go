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

	// var miners []*Miner
	// db.Where("mined_telegram > ?", Fee).Find(&miners)
	// log.Println(len(miners))

	// miners[1].TelegramId

	msg := `Please notice that you have anotes accrued on this bot!
	
To withdraw them, click here -> /withdraw

Or click the menu on this bot and choose withdraw. If you haven't attached the wallet (app.anotedao.com) yet, open it and click Connect Telegram on the bottom and do withdraw after that.

After withdrawal, you will receive your anotes once a day, when you enter daily mining code. Happy mining! 🚀`

	// telegramNotification(963770508, msg)

	mac.Run("127.0.0.1", Port)
}
