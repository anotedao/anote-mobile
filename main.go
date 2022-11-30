package main

import (
	"log"

	"gopkg.in/macaron.v1"
)

var conf *Config

var m *macaron.Macaron

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	m = initMacaron()

	// val := "%s%d%s%s__EVzrzuQUdTIPMtYlQaG+PL4uSnrkryR0DAY=__256857__m3e9byCWgCrQyxl4KGKvgvwXvSCW5VHOcXXKEw==__3A981xwTaRapfpk4cVWf1DN664ma1ztfd7a"
	// dataTransaction("3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p", &val, nil, nil)

	// db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }
	// ip := "167.99.36.116"
	// results, err := db.Get_all(ip)

	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }

	// log.Println(prettyPrint(results))

	m.Run("127.0.0.1", Port)
}
