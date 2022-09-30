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

	// height := int64(172000)
	// dataTransaction("3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT", nil, &height, nil)
	// dataTransaction("3AKGP29V8Pjh5VekzXq1SnwWXjMkQm7Zf9h", nil, nil, nil)

	// count := countReferred("3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p", nil)
	// log.Println(count)

	// addr := "3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT"
	// referral, _ := getData("referral", &addr)
	// log.Println(referral)

	// sendTelegramNotification("3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT")

	m.Run("127.0.0.1", Port)
}
