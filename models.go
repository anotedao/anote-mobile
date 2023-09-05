package main

import (
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type KeyValue struct {
	gorm.Model
	Key      string `gorm:"size:255;uniqueIndex"`
	ValueInt uint64 `gorm:"type:int"`
	ValueStr string `gorm:"type:string"`
}

type Miner struct {
	gorm.Model
	Address                string `gorm:"size:255;uniqueIndex"`
	LastNotification       time.Time
	LastNotificationWeekly time.Time `gorm:"default:'2023-06-17 23:00:00.797487649+00:00'"`
	TelegramId             int64     `gorm:"uniqueIndex"`
	MiningHeight           int64
	MiningTime             time.Time
	ReferralID             uint `gorm:"index"`
	Balance                uint64
	MinedTelegram          uint64
	MinedMobile            uint64
	LastPing               time.Time
	PingCount              int64
	IpAddresses            []*IpAddress `gorm:"many2many:miner_ip_addresses;"`
	UpdatedApp             bool         `gorm:"default:false"`
	LastInvite             time.Time
	BatteryNotification    bool `gorm:"default:false"`
	Cycles                 uint64
}

func (m *Miner) saveInBlockchain() {
	md := "%d%s__" + strconv.Itoa(int(m.MiningHeight))

	if m.ReferralID != 0 {
		r := &Miner{}
		db.First(r, m.ReferralID)
		md += "__" + r.Address
	}

	err := dataTransaction(m.Address, &md, nil, nil)

	counter := 0
	for err != nil && counter < 10 {
		time.Sleep(time.Millisecond * 500)
		err = dataTransaction(m.Address, &md, nil, nil)
		counter++
	}
}

func (m *Miner) saveIp(ip string) {
	ipa := &IpAddress{}
	db.FirstOrCreate(ipa, &IpAddress{Address: ip})
	db.Model(m).Association("IpAddresses").Append(ipa)
}

func (m *Miner) clearIps() {
	db.Model(m).Association("IpAddresses").Clear()
}

func getMiner(addr string) *Miner {
	mnr := &Miner{}
	db.First(mnr, &Miner{Address: addr})

	return mnr
}

func getMinerTel(tid int64) *Miner {
	mnr := &Miner{}
	db.First(mnr, &Miner{TelegramId: tid})

	return mnr
}

func getMinerOrCreate(addr string) *Miner {
	mnr := &Miner{}
	result := db.FirstOrCreate(mnr, &Miner{Address: addr})
	if result.Error != nil {
		log.Println(result.Error)
		logTelegram(result.Error.Error())
		log.Println(addr)
		return nil
	}

	return mnr
}

type IpAddress struct {
	gorm.Model
	Address string   `gorm:"size:255;uniqueIndex"`
	Miners  []*Miner `gorm:"many2many:miner_ip_addresses;"`
}
