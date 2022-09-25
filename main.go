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

	m.Run("127.0.0.1", Port)
}
