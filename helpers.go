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
)

func getMiningCode() int {
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
	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			// MaxConnsPerHost:   -1,
			MaxIdleConnsPerHost: -1,
			DisableKeepAlives:   true,
		},
	}})
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

	wc, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			// MaxConnsPerHost:   -1,
			MaxIdleConnsPerHost: -1,
			DisableKeepAlives:   true,
		},
	}})
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

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			// MaxConnsPerHost:   -1,
			MaxIdleConnsPerHost: -1,
			DisableKeepAlives:   true,
		},
	}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bh, _, err := cl.Blocks.Height(ctx)

	if err == nil {
		height = bh.Height
	}

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
	client, err := client.NewClient(client.Options{BaseUrl: nodeURL, Client: &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			// MaxConnsPerHost:   -1,
			MaxIdleConnsPerHost: -1,
			DisableKeepAlives:   true,
		},
	}})
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

func sendAssetTelegram(amount uint64, assetId string, recipient string) error {
	var networkByte byte
	var nodeURL string

	networkByte = 55
	nodeURL = AnoteNodeURL

	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKeyTelegram)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKeyTelegram)
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
	client, err := client.NewClient(client.Options{BaseUrl: nodeURL, Client: &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			// MaxConnsPerHost:   -1,
			MaxIdleConnsPerHost: -1,
			DisableKeepAlives:   true,
		},
	}})
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
	var amount int64
	var amountBasic int64
	var referralIndex float64
	miner := getMiner(address)
	stats := cch.StatsCache

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

		cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{
			Transport: &http.Transport{
				ForceAttemptHTTP2: true,
				// MaxConnsPerHost:   -1,
				MaxIdleConnsPerHost: -1,
				DisableKeepAlives:   true,
			},
		}})
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

		amount = (int64(total.Balance) / (int64(stats.ActiveUnits) + int64(stats.ActiveReferred/4))) - Fee

		if amount > 0 {
			amountBasic = amount

			rc := getRefCount(miner)

			if hasAintHealth(miner, true) {
				amount *= 10
			}

			referralIndex = float64(rc) * 0.25

			if heightDif > 2880 {
				times := int(heightDif / 1440)
				for i := 0; i < times; i++ {
					if amount > Fee {
						amount /= 2
					}
				}
				referralIndex = 1.0
			}

			fa := amount + int64(float64(amountBasic)*referralIndex)
			if fa > MULTI8 {
				log.Println(prettyPrint(total))
				log.Println(prettyPrint(stats))
				log.Println(fa)
				log.Println(amountBasic)
				log.Println(amount)
				log.Println(rc)
				logTelegram("Large amount issue.")
				fa = MULTI8
			}

			if strings.HasPrefix(address, "3A") {
				sendAsset(uint64(fa), "", address)
			}
		}
	}
}

func sendMinedTelegram(address string, heightDif int64) {
	var amount uint64
	var amountBasic uint64
	var referralIndex float64
	miner := getMiner(address)
	stats := cch.StatsCache

	if miner.ID != 0 {
		amount = (4320000000 / (uint64(stats.ActiveUnits) + uint64(stats.ActiveReferred/4))) - Fee
		amountBasic = amount

		rc := getRefCount(miner)

		referralIndex = float64(rc) * 0.25

		if heightDif > 2880 {
			times := int(heightDif / 1440)
			for i := 0; i < times; i++ {
				if amount > Fee {
					amount /= 2
				}
			}
			referralIndex = 1.0
		}

		fa := amount + uint64(float64(amountBasic)*referralIndex)
		if fa > MULTI8 {
			fa = MULTI8
		}

		if !isFollower(miner.TelegramId) {
			fa = uint64(float64(fa) * 0.9)
		}

		if strings.HasPrefix(address, "3A") {
			sendAssetTelegram(fa, "", address)
		}
	}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func sendTelegramNotification(addr string, height int64, savedHeight int64) bool {
	resp, err := http.Get(fmt.Sprintf("http://localhost:5006/notification/%s/%d/%d", addr, height, savedHeight))
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

func telegramNotification(tid int64, msg string) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:5006/notification-tg/%d/%s", tid, msg))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
	defer resp.Body.Close()
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

	resp, err := http.Get(fmt.Sprintf("http://localhost:5006/log/%s", message))
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
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

	for _, m := range miners {
		if height-uint64(m.MiningHeight) <= 1440 {
			sr.ActiveMiners++
			if m.ReferralID != 0 {
				sr.ActiveReferred++
			}
		}

		if height-uint64(m.MiningHeight) <= 1440 {
			sr.PayoutMiners++

			if hasAintHealth(m, true) {
				sr.ActiveUnits += 10
			} else {
				sr.ActiveUnits++
			}
		}
	}

	sr.InactiveMiners = len(miners) - sr.PayoutMiners

	return sr
}

type Stats struct {
	ActiveMiners   int `json:"active_miners"`
	ActiveReferred int `json:"active_referred"`
	PayoutMiners   int `json:"payout_miners"`
	InactiveMiners int `json:"inactive_miners"`
	ActiveUnits    int `json:"active_units"`
}

func getRefCount(m *Miner) uint64 {
	var miners []*Miner

	height := getHeight()

	db.Where("referral_id = ? AND mining_height > ?", m.ID, height-2880).Find(&miners)
	count := len(miners)

	miners = nil

	return uint64(count)
}

func hasAintHealth(m *Miner, second bool) bool {
	sma := StakeMobileAddress

	d, err := getData("%s__"+m.Address, &sma)
	if err != nil || d == nil {
		return false
	}

	aint := parseItem(d.(string), 0)
	if aint != nil {
		if second && aint.(int) >= (10*MULTI8) {
			return true
		} else if !second && aint.(int) >= MULTI8 {
			return true
		}
	}

	return false
}

func sendInvite(m *Miner) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:5006/invite/%s", strconv.Itoa(int(m.TelegramId))))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
	defer resp.Body.Close()
}

func sendNotificationEnd(m *Miner) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:5006/notification-end/%s", strconv.Itoa(int(m.TelegramId))))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
	defer resp.Body.Close()
}

func sendNotificationWeekly(m *Miner) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:5006/notification-weekly/%s", strconv.Itoa(int(m.TelegramId))))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
	defer resp.Body.Close()
}

func sendNotificationBattery(m *Miner) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:5006/notification-bo/%s", strconv.Itoa(int(m.TelegramId))))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
	defer resp.Body.Close()
}

func sendNotificationFirst(m *Miner) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:5006/notification-first/%s", strconv.Itoa(int(m.TelegramId))))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
	defer resp.Body.Close()
}

func getBalance(address string) (uint64, error) {
	addr, err := proto.NewAddressFromString(address)
	if err != nil {
		return 0, err
	}

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			// MaxConnsPerHost:   -1,
			MaxIdleConnsPerHost: -1,
			DisableKeepAlives:   true,
		},
	}})
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	total, _, err := cl.Addresses.Balance(ctx, addr)
	if err != nil {
		return 0, err
	}

	return total.Balance, nil
}

func getMiningFactor(m *Miner) float64 {
	mf := float64(1)

	rc := getRefCount(m)

	referralIndex := float64(rc) * 0.25

	if hasAintHealth(m, true) {
		mf *= 10
	}

	mf += referralIndex

	return mf
}

func getBasicAmount(amount uint64) uint64 {
	ba := uint64(0)

	if cch != nil && cch.StatsCache != nil {
		stats := cch.StatsCache
		ba = uint64(float64(amount) / float64((float64(stats.ActiveUnits) + float64(stats.ActiveReferred)/4)))
	}

	// float64((float64(totalt.Balance) / float64(uint64(stats.ActiveUnits)+uint64(stats.ActiveReferred/4)))) / MULTI8

	return ba
}

type AlphaSentResponse struct {
	Sent bool `json:"sent"`
}

func getAlphaSent(addr string) bool {
	alr := &AlphaSentResponse{Sent: true}
	resp, err := http.Get(fmt.Sprintf("http://localhost:5006/alpha-sent/%s", addr))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return true
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, alr); err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return true
	}

	return alr.Sent
}

func isFollower(tid int64) bool {
	ifr := &IsFollowerResponse{IsFollower: false}
	resp, err := http.Get(fmt.Sprintf("http://localhost:5006/is-follower/%d", tid))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return true
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, ifr); err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return true
	}

	return ifr.IsFollower
}

type IsFollowerResponse struct {
	IsFollower bool `json:"is_follower"`
}
