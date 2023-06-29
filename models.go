package main

import (
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
	Address                string `gorm:"size:255"`
	LastNotification       time.Time
	LastNotificationWeekly time.Time `gorm:"default:'2023-06-17 23:00:00.797487649+00:00'"`
	TelegramId             int64     `gorm:"uniqueIndex"`
	MiningHeight           int64
	MiningTime             time.Time
	ReferralID             uint `gorm:"index"`
	Confirmed              bool `gorm:"default:false"`
	Balance                uint64
	MinedTelegram          uint64
	MinedMobile            uint64
	LastPing               time.Time
	PingCount              int64
	IpAddresses            []*IpAddress `gorm:"many2many:miner_ip_addresses;"`
	UpdatedApp             bool         `gorm:"default:false"`
	LastInvite             time.Time
	BatteryNotification    bool `gorm:"default:false"`
}

func (m *Miner) saveInBlockchain() {
	md := "%d%s__" + strconv.Itoa(int(m.MiningHeight))

	if m.ReferralID != 0 {
		r := &Miner{}
		db.First(r, m.ReferralID)
		md += "__" + r.Address
	}

	err := dataTransaction(m.Address, &md, nil, nil)

	go func(md string, err error) {
		for err != nil {
			time.Sleep(time.Millisecond * 500)
			err = dataTransaction(m.Address, &md, nil, nil)
		}
	}(md, err)
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
	m := &Miner{}
	db.First(m, &Miner{Address: addr})

	return m
}

func getMinerTel(tid int64) *Miner {
	m := &Miner{}
	db.FirstOrCreate(m, &Miner{TelegramId: tid})

	return m
}

func getMinerOrCreate(addr string) *Miner {
	m := &Miner{}
	result := db.FirstOrCreate(m, &Miner{Address: addr})
	if result.RowsAffected == 1 {
		mon.Miners = append(mon.Miners, m)
	}

	return m
}

type IpAddress struct {
	gorm.Model
	Address string   `gorm:"size:255;uniqueIndex"`
	Miners  []*Miner `gorm:"many2many:miner_ip_addresses;"`
}
