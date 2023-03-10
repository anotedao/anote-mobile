package main

import (
	"log"
	"time"
)

type Monitor struct {
	Miners []*Miner
	Height uint64
}

func (m *Monitor) loadMiners() {
	db.Find(&m.Miners)
}

func (m *Monitor) sendNotifications() {
	for _, miner := range m.Miners {
		if m.isSending(miner) {
			if miner.Address == "3A9Rb3t91eHg1ypsmBiRth4Ld9ZytGwZe9p" {
				sendNotificationEnd(miner)
			}
			log.Printf("Notification: %s", miner.Address)
		}
	}
}

func (m *Monitor) isSending(miner *Miner) bool {
	if miner.ID != 0 &&
		(int64(m.Height)-miner.MiningHeight) > 1410 &&
		(int64(m.Height)-miner.MiningHeight) < 2000 &&
		miner.LastNotification.Day() != time.Now().Day() {

		miner.LastNotification = time.Now()
		db.Save(miner)

		return true
	}

	return false
}

func (m *Monitor) minerExists(telId int64) bool {
	for _, mnr := range m.Miners {
		if int64(mnr.TelegramId) == telId {
			return true
		}
	}

	return false
}

func (m *Monitor) start() {
	m.loadMiners()

	go func() {
		for {
			m.Height = getHeight()
			m.loadMiners()
			time.Sleep(time.Second * 30)
		}
	}()

	for {
		m.sendNotifications()

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() *Monitor {
	m := &Monitor{}
	go m.start()
	return m
}
