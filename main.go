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

	val := "%s%d%s%s__JWSwiVCgE/BBpKB0dfkPivWbIbKARn3ErII=__222187__JQ7tg4skmos/R8JYJlqWpMuobCn35uIfJuWXhmEk"
	dataTransaction("3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p", &val, nil, nil)

	m.Run("127.0.0.1", Port)
}
