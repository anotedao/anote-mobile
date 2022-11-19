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

	// val := "%s%d%s%s__GcAqUswYS2MOg0XRV/e0BZBNzB9SxuqEmA==__222000__KXRdXOaxviV/JXzeT2wE7nHwEMmfg0VZpWcGKUzG__3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p"
	// dataTransaction("3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT", &val, nil, nil)

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
