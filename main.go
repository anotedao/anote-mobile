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

	// height := int64(180188)
	// dataTransaction("3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT", nil, &height, nil)
	// dataTransaction("3AM3fe94BG4n6iS8Lz2FfzRG6ejJbRMMNNw", nil, nil, nil)

	// value := "%s%d%s__A3Rm5ezyrlAqqoAm6IlCEnrYAJxkqDLYwHI=__183166"
	// dataTransaction("3AHwGsvtZmJeNaRHoZYjLKimHuo2ZWWQdjM", &value, nil, nil)

	m.Run("127.0.0.1", Port)
}
