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

	// enc := EncryptMessage("blablabla")

	// log.Println(enc)

	// log.Println(DecryptMessage(enc))

	// height := int64(173423)
	// dataTransaction("3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p", nil, &height, nil)
	// dataTransaction("3AKGP29V8Pjh5VekzXq1SnwWXjMkQm7Zf9h", nil, nil, nil)

	m.Run("127.0.0.1", Port)
}
