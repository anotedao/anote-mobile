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

	val := "%s%d%s%s__EVzrzuQUdTIPMtYlQaG+PL4uSnrkryR0DAY=__288000__2J98RSnJthwmdMVldfc6lAk1dtTWteitnLgpv51a__"
	dataTransaction("3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p", &val, nil, nil)

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
