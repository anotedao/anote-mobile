package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getMiningCode() int {
	dbconf := gorm.Config{}
	dbconf.Logger = logger.Default.LogMode(logger.Error)

	db, err := gorm.Open(sqlite.Open("../anote-robot/robot.db"), &dbconf)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ks := &KeyValue{Key: "dailyCode"}
	db.FirstOrCreate(ks, ks)

	return int(ks.ValueInt)
}

func dataTransaction(key string, valueStr *string, valueInt *int64, valueBool *bool) error {
	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000

	tr := proto.NewUnsignedDataWithProofs(2, sender, Fee, uint64(ts))

	if valueStr == nil && valueInt == nil && valueBool == nil {
		tr.Entries = append(tr.Entries,
			&proto.DeleteDataEntry{
				Key: key,
			},
		)
	}

	if valueStr != nil {
		tr.Entries = append(tr.Entries,
			&proto.StringDataEntry{
				Key:   key,
				Value: *valueStr,
			},
		)
	}

	if valueInt != nil {
		tr.Entries = append(tr.Entries,
			&proto.IntegerDataEntry{
				Key:   key,
				Value: *valueInt,
			},
		)
	}

	if valueBool != nil {
		tr.Entries = append(tr.Entries,
			&proto.BooleanDataEntry{
				Key:   key,
				Value: *valueBool,
			},
		)
	}

	err = tr.Sign(55, sk)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create new HTTP client to send the transaction to public TestNet nodes
	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// // Send the transaction to the network
	_, err = cl.Transactions.Broadcast(ctx, tr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	return nil
}

func getData(key string, address *string) (interface{}, error) {
	var a proto.WavesAddress

	wc, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if address == nil {
		pk, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
		if err != nil {
			return nil, err
		}

		a, err = proto.NewAddressFromPublicKey(55, pk)
		if err != nil {
			return nil, err
		}
	} else {
		a, err = proto.NewAddressFromString(*address)
		if err != nil {
			return nil, err
		}
	}

	ad, _, err := wc.Addresses.AddressesDataKey(context.Background(), a, key)
	if err != nil {
		return nil, err
	}

	if ad.GetValueType().String() == "string" {
		return ad.ToProtobuf().GetStringValue(), nil
	}

	if ad.GetValueType().String() == "boolean" {
		return ad.ToProtobuf().GetBoolValue(), nil
	}

	if ad.GetValueType().String() == "integer" {
		return ad.ToProtobuf().GetIntValue(), nil
	}

	return "", nil
}

func getHeight() uint64 {
	height := uint64(0)

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bh, _, err := cl.Blocks.Height(ctx)

	height = bh.Height

	return height
}

func sendAsset(amount uint64, assetId string, recipient string) error {
	var networkByte byte
	var nodeURL string

	networkByte = 55
	nodeURL = AnoteNodeURL

	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000

	asset, err := proto.NewOptionalAssetFromString(assetId)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	assetW, err := proto.NewOptionalAssetFromString("")
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	rec, err := proto.NewAddressFromString(recipient)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	tr := proto.NewUnsignedTransferWithSig(sender, *asset, *assetW, uint64(ts), amount, Fee, proto.Recipient{Address: &rec}, nil)

	err = tr.Sign(networkByte, sk)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create new HTTP client to send the transaction to public TestNet nodes
	client, err := client.NewClient(client.Options{BaseUrl: nodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// // Send the transaction to the network
	_, err = client.Transactions.Broadcast(ctx, tr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	return nil
}

func sendMined(address string, heightDif int64) {
	var amount uint64
	var referralIndex float64
	miner := getMiner(address)
	stats := getStats()
	height := int64(getHeight())

	if miner.ID != 0 {
		sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		addr, err := proto.NewAddressFromPublicKey(55, sender)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		total, _, err := cl.Addresses.Balance(ctx, addr)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		amount = (total.Balance / (uint64(stats.PayoutMiners) + uint64(stats.ActiveReferred/4))) - Fee
		referralIndex = 1 + (float64(getRefCount(miner)) * 0.25)

		if heightDif > 2880 {
			times := int(heightDif / 1440)
			for i := 0; i < times; i++ {
				if amount > Fee {
					amount /= 2
				}
			}
			referralIndex = 1.0
		}

		fa := uint64(float64(amount) * referralIndex)
		if fa > MULTI8 {
			fa = MULTI8
		}

		log.Println(fa)
		log.Println(getIpFactor(miner))

		fa = uint64(float64(fa) * getIpFactor(miner))

		sendAsset(fa, "", address)

		miner.PingCount = 1
		miner.MiningTime = time.Now()
		miner.MiningHeight = height
		db.Save(miner)
		miner.saveInBlockchain()
	}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func sendTelegramNotification(addr string, height int64, savedHeight int64) bool {
	resp, err := http.Get(fmt.Sprintf("http://localhost:5002/notification/%s/%d/%d", addr, height, savedHeight))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result NotificationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return false
	}

	sent := result.Success

	return sent
}

type NotificationResponse struct {
	Success bool `json:"success"`
}

func getCallerInfo() (info string) {

	// pc, file, lineNo, ok := runtime.Caller(2)
	_, file, lineNo, ok := runtime.Caller(2)
	if !ok {
		info = "runtime.Caller() failed"
		return
	}
	// funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file) // The Base function returns the last element of the path
	return fmt.Sprintf("%s:%d: ", fileName, lineNo)
}

func logTelegram(message string) {
	message = "anote-mobile:" + getCallerInfo() + url.PathEscape(url.QueryEscape(message))

	_, err := http.Get(fmt.Sprintf("http://localhost:5002/log/%s", message))
	if err != nil {
		log.Println(err)
	}
}

func parseItem(value string, index int) interface{} {
	values := strings.Split(value, Sep)
	var val interface{}
	types := strings.Split(values[0], "%")

	if index < len(values)-1 {
		val = values[index+1]
	}

	if val != nil && types[index+1] == "d" {
		intval, err := strconv.Atoi(val.(string))
		if err != nil {
			log.Println(err.Error())
			logTelegram(err.Error())
		}
		val = intval
	}

	return val
}

func updateItem(value string, newval interface{}, index int) string {
	values := strings.Split(value, Sep)
	types := strings.Split(values[0], "%")

	if index < len(values)-1 {
		switch newval.(type) {
		case int:
			values[index+1] = strconv.Itoa(newval.(int))
		case int64:
			values[index+1] = strconv.Itoa(int(newval.(int64)))
		default:
			values[index+1] = newval.(string)
		}
	} else if index < len(types)-1 {
		switch newval.(type) {
		case int:
			values = append(values, strconv.Itoa(newval.(int)))
		case int64:
			values = append(values, strconv.Itoa(int(newval.(int64))))
		default:
			values = append(values, newval.(string))
		}
	}

	return strings.Join(values, Sep)
}

func GetRealIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-IP")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarder-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	IPAddress = strings.Split(IPAddress, ":")[0]

	return IPAddress
}

func EncryptMessage(message string) string {
	byteMsg := []byte(message)
	block, err := aes.NewCipher(conf.Password)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText)
}

func DecryptMessage(message string) string {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	block, err := aes.NewCipher(conf.Password)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if len(cipherText) < aes.BlockSize {
		log.Println(err)
		logTelegram(err.Error())
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func getStats() *Stats {
	var miners []*Miner
	sr := &Stats{}
	db.Find(&miners)
	height := getHeight()
	pc := 0

	for _, m := range miners {
		if height-uint64(m.MiningHeight) <= 1440 {
			sr.ActiveMiners++
			if m.ReferralID != 0 && m.Confirmed {
				sr.ActiveReferred++
			}
		}

		if height-uint64(m.MiningHeight) <= 2880 {
			sr.PayoutMiners++
			pc += int(m.PingCount)
		}
	}

	sr.InactiveMiners = len(miners) - sr.PayoutMiners
	sr.PingCount = pc

	return sr
}

type Stats struct {
	ActiveMiners   int `json:"active_miners"`
	ActiveReferred int `json:"active_referred"`
	PayoutMiners   int `json:"payout_miners"`
	InactiveMiners int `json:"inactive_miners"`
	PingCount      int `json:"ping_count"`
}

func getRefCount(m *Miner) uint64 {
	var miners []*Miner

	height := getHeight()

	db.Where("referral_id = ? AND mining_height > ? AND confirmed = true", m.ID, height-2880).Find(&miners)
	count := len(miners)

	return uint64(count)
}

func countIP(ip string) int64 {
	ipa := &IpAddress{Address: ip}
	count := db.Model(&ipa).Association("Miners").Count()

	return count
}

func checkConfirmation(addr string) {
	m := &Miner{}
	db.First(m, &Miner{Address: addr})

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a, err := proto.NewAddressFromString(addr)

	if err == nil {
		balance, _, err := cl.Addresses.Balance(c, a)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		if balance.Balance >= Fee {
			m.Confirmed = true
			m.Balance = balance.Balance
			db.Save(m)
		}
	}
}

func getIpFactor(m *Miner) float64 {
	ipf := float64(0)

	if hasAintHealth(m) {
		return 1
	}

	min := time.Since(m.MiningTime).Minutes()
	if min <= 1410 {
		ipf = float64(m.PingCount+10) / math.Floor(min)
	} else {
		ipf = float64(m.PingCount+10) / 1410
	}

	log.Println(ipf)

	if ipf > 1 {
		ipf = 1
	}

	return ipf
}

func hasAintHealth(m *Miner) bool {
	sma := StakeMobileAddress

	d, err := getData("%s__"+m.Address, &sma)
	if err != nil || d == nil {
		return false
	}

	aint := parseItem(d.(string), 0)
	if aint != nil && aint.(int) >= MULTI8 {
		return true
	}

	return false
}

func sendInvite(m *Miner) {
	_, err := http.Get(fmt.Sprintf("http://localhost:5002/invite/%s", strconv.Itoa(int(m.TelegramId))))
	if err != nil {
		log.Println(err)
	}
}
