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

	// cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	// if err != nil {
	// 	log.Println(err)
	// }

	// // Context to cancel the request execution on timeout
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
	// if err != nil {
	// 	log.Println(err)
	// }

	// addr, err := proto.NewAddressFromPublicKey(55, sender)
	// if err != nil {
	// 	log.Println(err)
	// }

	// data, _, err := cl.Addresses.AddressesData(ctx, addr)

	// for _, de := range data {
	// 	val := de.ToProtobuf().GetStringValue()
	// 	vala := strings.Split(val, Sep)
	// 	valNew := "%s%d%s%s" + Sep + vala[1]
	// 	if len(vala) > 2 {
	// 		valNew += Sep + vala[2] + Sep + EncryptMessage("127.0.0.1")
	// 	}
	// 	if len(vala) == 4 {
	// 		valNew += Sep + vala[3]
	// 	}
	// 	log.Println(valNew)
	// 	dataTransaction(de.GetKey(), &valNew, nil, nil)
	// }

	log.Println(DecryptMessage("FyikphXCxu5QIHhbFkVoq+JCYt2xOxfPvGNmCOVbBg=="))

	m.Run("127.0.0.1", Port)
}
