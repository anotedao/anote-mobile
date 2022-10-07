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
	// dataTransaction("3AJ4g8UdJuDiSrE9BhmSuSeknK265hG5XA1", nil, nil, nil)

	// value := "%s%d%s__aR22W2epPnoD+cx+OpX0aN4I24kJK+WelQ==__187423__3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p"
	// dataTransaction("3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT", &value, nil, nil)

	m.Run("127.0.0.1", Port)
}
