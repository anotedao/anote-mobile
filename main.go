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

	// height := int64(180188)
	// dataTransaction("3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT", nil, &height, nil)
	// dataTransaction("3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p", nil, &height, nil)
	// dataTransaction("3ASTzMMBPV6GXWxQkGtYdsXciaTGG92Yixo", nil, nil, nil)
	// dataTransaction("3AF4JjMnExbNYYxDF3AKes6Ce1M1NyuSYz7", nil, nil, nil)
	// dataTransaction("3AKCefhcrijSwwWM671ahhMrPVrE7Je3j4s", nil, nil, nil)
	// dataTransaction("3ASLefwuE3dz9cW9bhP6ZC9N73pLqs2vPEH", nil, nil, nil)

	// count := countReferred("3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p", nil)
	// log.Println(count)

	// addr := "3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT"
	// referral, _ := getData("referral", &addr)
	// log.Println(referral)

	// sendTelegramNotification("3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT")

	// value, _ := getData("3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p", nil)
	// newval := updateItem(value.(string), 12, 0)
	// log.Println(newval)

	m.Run("127.0.0.1", Port)
}
