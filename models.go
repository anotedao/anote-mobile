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
	Address          string `gorm:"size:255;uniqueIndex"`
	LastNotification time.Time
	TelegramId       int64
	MiningHeight     int64
	MiningTime       time.Time
	ReferralID       uint   `gorm:"index"`
	IP               string `gorm:"index"`
	Confirmed        bool   `gorm:"default:false"`
	Balance          uint64
	LastPing         time.Time
	PingCount        int64
	IP2              string `gorm:"index"`
	IP3              string `gorm:"index"`
	IP4              string `gorm:"index"`
	IP5              string `gorm:"index"`
}

func (m *Miner) saveInBlockchain() {
	md := "%d%s__" + strconv.Itoa(int(m.MiningHeight))

	if m.ReferralID != 0 {
		r := &Miner{}
		db.First(r, m.ReferralID)
		md += "__" + r.Address
	}

	dataTransaction(m.Address, &md, nil, nil)
}

func getMiner(addr string) *Miner {
	m := &Miner{}
	db.First(m, &Miner{Address: addr})

	return m
}

func getMinerOrCreate(addr string) *Miner {
	m := &Miner{}
	db.FirstOrCreate(m, &Miner{Address: addr})

	return m
}
