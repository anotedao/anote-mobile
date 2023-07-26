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
	counter := 1
	for _, miner := range m.Miners {
		if m.isSending(miner, 1410) {
			sendNotificationEnd(miner)
			log.Printf("Notification: %s", miner.Address)
		}

		if m.isSendingWeekly(miner, 10080) {
			sendNotificationWeekly(miner)
			log.Printf("Notification Weekly: %s %d", miner.Address, counter)
			counter++
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
		err := db.Save(miner).Error
		for err != nil {
			time.Sleep(time.Millisecond * 500)
			err = db.Save(miner).Error
			log.Println(err)
		}

		return true
	}

	return false
}

func (m *Monitor) isSendingWeekly(miner *Miner, limit int64) bool {
	if miner.ID != 0 &&
		(int64(m.Height)-miner.MiningHeight) >= limit &&
		(miner.MiningTime.Hour() == time.Now().Hour() ||
			miner.MiningTime.IsZero()) &&
		time.Since(miner.LastNotificationWeekly) > time.Hour*168 &&
		miner.TelegramId != 0 {

		miner.LastNotificationWeekly = time.Now()
		err := db.Save(miner).Error
		for err != nil {
			time.Sleep(time.Millisecond * 500)
			err = db.Save(miner).Error
			log.Println(err)
		}

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

			log.Printf("Telegram Basic: %d", ba)

			m.OldBalanceTelegram = m.NewBalanceTelegram

			ks := &KeyValue{Key: "oldBalanceTelegram"}
			db.FirstOrCreate(ks, ks)
			ks.ValueInt = m.OldBalanceTelegram
			err := db.Save(ks).Error
			for err != nil {
				time.Sleep(time.Millisecond * 500)
				err = db.Save(ks).Error
				log.Println(err)
			}

			for _, mnr := range m.Miners {
				if m.Height-uint64(mnr.MiningHeight) <= 1410 {
					mnr.MinedTelegram += uint64(float64(ba) * getMiningFactor(mnr))
					err := db.Save(mnr).Error
					for err != nil {
						time.Sleep(time.Millisecond * 500)
						err = db.Save(mnr).Error
						log.Println(err)
					}
				}
			}

			log.Printf("New Telegram Amount: %d", diff)
		}
	} else if err != nil {
		log.Println(err)
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
			time.Sleep(time.Second * 120)
		}
	}()

	// go func() {
	// 	for {
	// 		m.loadMiners()
	// 		time.Sleep(time.Minute * 2)
	// 	}
	// }()

	total := 0

	for _, mnr := range m.Miners {
		total += int(mnr.MinedTelegram)
	}

	log.Printf("Total Telegram: %d", total)

	for {
		m.loadMiners()

		m.sendNotifications()

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() *Monitor {
	m := &Monitor{}
	go m.start()
	return m
}
