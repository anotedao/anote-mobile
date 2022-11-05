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

	val := "%s%d%s%s__GcAqUswYS2MOg0XRV/e0BZBNzB9SxuqEmA==__222000__KXRdXOaxviV/JXzeT2wE7nHwEMmfg0VZpWcGKUzG__3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p"
	dataTransaction("3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT", &val, nil, nil)

	m.Run("127.0.0.1", Port)
}
