package main

import (
	"log"

	"gopkg.in/macaron.v1"
	"gorm.io/gorm"
)

var conf *Config

var m *macaron.Macaron

var db *gorm.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	db = initDb()

	m = initMacaron()

	// val := "%d%s__31000__3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p"
	// dataTransaction("3AShXVgRcRis82CwD7o9pz1Ac9vmRYMqELT", &val, nil, nil)

	m.Run("127.0.0.1", Port)

	// type DataItems []struct {
	// 	Key   string `json:"key"`
	// 	Type  string `json:"type"`
	// 	Value string `json:"value"`
	// }

	// di := &DataItems{}

	// file, err := os.Open("data.json")

	// if err != nil {
	// 	log.Println(err)
	// }

	// decoder := json.NewDecoder(file)

	// err = decoder.Decode(di)

	// if err != nil {
	// 	log.Println(err)
	// }

	// count := 0

	// for i, d := range *di {
	// 	r := parseItem(d.Value, 3)

	// 	if r != nil && strings.HasPrefix(r.(string), "3A") {
	// 		count++
	// 	}
	// 	log.Println(i)
	// }
	// log.Println(count)

	// var miners []*Miner
	// db.Find(&miners)

	// for i, m := range miners {
	// 	val := "%d%s__0"

	// 	if m.ReferralID != 0 {
	// 		r := &Miner{}
	// 		db.First(r, m.ReferralID)
	// 		val += "__" + r.Address
	// 	}

	// 	dataTransaction(m.Address, &val, nil, nil)
	// 	log.Println(i)
	// }
}
