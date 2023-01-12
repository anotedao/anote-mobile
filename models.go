package main

import (
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
	ReferralID       uint   `gorm:"index"`
	IP               string `gorm:"index;default:127.0.0.1"`
	Confirmed        bool   `gorm:"default:false"`
	Balance          uint64
	LastPing         time.Time
	PingCount        int64
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
