package main

import (
	"log"
	"time"
)

type Monitor struct {
	Miners             []*Miner
	Height             uint64
	OldBalanceTelegram uint64
	NewBalanceTelegram uint64
}

func (m *Monitor) loadMiners() {
	db.Find(&m.Miners)
}

func (m *Monitor) sendNotifications() {
	for _, miner := range m.Miners {
		if m.isSending(miner, 1410) {
			sendNotificationEnd(miner)
			log.Printf("Notification: %s", miner.Address)
		}

		// if m.isSendingBattery(miner) {
		// 	sendNotificationBattery(miner)
		// 	log.Printf("Notification battery: %s", miner.Address)
		// }
	}
}

func (m *Monitor) isSending(miner *Miner, limit int64) bool {
	if miner.ID != 0 &&
		(int64(m.Height)-miner.MiningHeight) >= limit &&
		(int64(m.Height)-miner.MiningHeight) < limit+30 &&
		miner.LastNotification.Day() != time.Now().Day() &&
		miner.MiningHeight > 0 &&
		miner.TelegramId != 0 {

		miner.LastNotification = time.Now()
		db.Save(miner)

		return true
	}

	return false
}

func (m *Monitor) isSendingBattery(miner *Miner) bool {
	// health := int(getIpFactor(miner, true, uint64(m.Height), 2) * 100)

	// if health > 100 {
	// 	health = 100
	// } else if health < 0 {
	// 	health = 0
	// }

	// if miner.ID != 0 &&
	// 	miner.BatteryNotification &&
	// 	health < 100 &&
	// 	time.Since(miner.MiningTime) > (time.Minute*5) &&
	// 	miner.TelegramId != 0 {

	// 	miner.BatteryNotification = false
	// 	db.Save(miner)

	// 	return true
	// }

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

func (m *Monitor) checkMined() {
	var err error
	m.NewBalanceTelegram, err = getBalance(TelegramAddress)
	if err == nil && m.NewBalanceTelegram > m.OldBalanceTelegram {
		diff := m.NewBalanceTelegram - m.OldBalanceTelegram
		if diff > 0 {
			ba := getBasicAmount(diff)

			m.OldBalanceTelegram = m.NewBalanceTelegram

			ks := &KeyValue{Key: "oldBalanceTelegram"}
			db.FirstOrCreate(ks, ks)
			ks.ValueInt = m.OldBalanceTelegram
			db.Save(ks)

			for _, mnr := range m.Miners {
				mnr.MinedTelegram += uint64(float64(ba) * getMiningFactor(mnr))
				db.Save(mnr)
			}

			log.Printf("New Telegram Amount: %d", diff)
		}
	}
}

func (m *Monitor) start() {
	m.loadMiners()

	ks := &KeyValue{Key: "oldBalanceTelegram"}
	db.FirstOrCreate(ks, ks)

	m.OldBalanceTelegram = uint64(ks.ValueInt)

	go func() {
		for {
			m.Height = getHeight()
			m.checkMined()
			time.Sleep(time.Second * 30)
		}
	}()

	go func() {
		for {
			m.loadMiners()
			time.Sleep(time.Minute * 2)
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
