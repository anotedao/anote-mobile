package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	macaron "gopkg.in/macaron.v1"
)

// func mineView(ctx *macaron.Context, cpt *captcha.Captcha) {
// 	height := int64(getHeight())

// 	pr := &MineResponse{
// 		Success: true,
// 		Error:   0,
// 	}

// 	addr := ctx.Params("address")
// 	if len(addr) > 0 && addr != "0" {
// 		cpid := ctx.Params("captchaid")
// 		cp := ctx.Params("captcha")
// 		code := ctx.Params("code")
// 		ip := GetRealIP(ctx.Req.Request)

// 		miner := getMinerOrCreate(addr)
// 		if miner != nil {
// 			savedHeight := miner.MiningHeight

// 			code = strings.TrimSpace(code)
// 			code = regexp.MustCompile(`[^0-9]+`).ReplaceAllString(code, "")

// 			codeInt, err := strconv.Atoi(code)
// 			if err != nil {
// 				log.Println(err)
// 				logTelegram(err.Error())
// 				pr.Success = false
// 				pr.Error = 2
// 			}

// 			if !cpt.Verify(cpid, cp) {
// 				pr.Success = false
// 				pr.Error = 1
// 			}

// 			if int(codeInt) != getMiningCode() {
// 				pr.Success = false
// 				pr.Error = 2
// 			}

// 			if !strings.HasPrefix(addr, "3A") {
// 				pr.Success = false
// 				pr.Error = 5
// 			}

// 			if pr.Error == 0 && (height-miner.MiningHeight > 1410) {
// 				log.Println(fmt.Sprintf("%s %s", addr, ip))

// 				miner.clearIps()
// 				miner.saveIp(ip)

// 				if savedHeight > 0 {
// 					sendMined(addr, height-savedHeight)
// 					sendMinedTelegram(addr, height-savedHeight)

// 					miner.Cycles++
// 					miner.MiningTime = time.Now()
// 					miner.MiningHeight = height
// 					miner.BatteryNotification = true
// 					err = db.Save(miner).Error
// 					if err != nil {
// 						log.Println(err)
// 						logTelegram(err.Error())
// 					}
// 					miner.saveInBlockchain()
// 				} else {
// 					sendAssetTelegram(Fee, "", addr)
// 					miner.MinedTelegram = Fee
// 					miner.Cycles = 1
// 					miner.MiningTime = time.Now()
// 					miner.MiningHeight = height
// 					miner.UpdatedApp = true
// 					miner.BatteryNotification = true
// 					if miner.Address == "" {
// 						miner.Address = strconv.Itoa(int(miner.TelegramId))
// 					}
// 					err := db.Save(miner).Error
// 					if err != nil {
// 						log.Println(err)
// 						logTelegram(err.Error())
// 					}
// 					miner.saveInBlockchain()
// 					sendNotificationFirst(miner)
// 				}
// 				// mon.loadMiners()
// 			}
// 		} else {
// 			pr.Success = false
// 			pr.Error = 7
// 		}
// 	} else {
// 		pr.Success = false
// 		pr.Error = 6
// 	}

// 	ctx.Resp.Header().Add("Access-Control-Allow-Origin", "*")
// 	ctx.JSON(200, pr)
// }

// func newCaptchaView(ctx *macaron.Context, cpt *captcha.Captcha) {
// 	c, err := cpt.CreateCaptcha()
// 	if err != nil {
// 		log.Println(err)
// 		logTelegram(err.Error())
// 	}

// 	ir := &ImageResponse{
// 		Id:    c,
// 		Image: fmt.Sprintf("%s/captcha/%s.png", conf.Host, c),
// 	}

// 	ctx.Resp.Header().Add("Access-Control-Allow-Origin", "*")
// 	ctx.JSON(200, ir)
// }

type MineResponse struct {
	Success bool `json:"success"`
	Error   int  `json:"error"`
}

type MinerResponse struct {
	ID            uint    `json:"id"`
	Address       string  `json:"address"`
	Referred      int64   `json:"referred"`
	Active        int64   `json:"active"`
	HasTelegram   bool    `json:"has_telegram"`
	MiningHeight  int64   `json:"mining_height"`
	Height        uint64  `json:"height"`
	Exists        bool    `json:"exists"`
	MinedMobile   uint64  `json:"mined_mobile"`
	MinedTelegram uint64  `json:"mined_telegram"`
	TelegramId    int64   `json:"telegram_id"`
	AlphaSent     bool    `json:"alpha_sent"`
	Cycles        uint64  `json:"cycles"`
	Price         float64 `json:"price"`
}

type HealthResponse struct {
	Health     int  `json:"health"`
	UpdatedApp bool `json:"updated_app"`
}

type ImageResponse struct {
	Image string `json:"image"`
	Id    string `json:"id"`
}

// func healthView(ctx *macaron.Context) {
// 	a := ctx.Params("address")

// 	hr := &HealthResponse{}

// 	miner := getMiner(a)

// 	hr.Health = 100

// 	if hr.Health > 100 {
// 		hr.Health = 100
// 	} else if hr.Health < 0 {
// 		hr.Health = 0
// 	}

// 	hr.UpdatedApp = miner.UpdatedApp

// 	ctx.JSON(200, hr)
// }

func statsView(ctx *macaron.Context) {
	sr := cch.StatsCache
	ctx.JSON(200, sr)
}

// func newUserView(ctx *macaron.Context) {
// 	u := &Miner{}
// 	r := &Miner{}

// 	ap := ctx.Params("addr")
// 	rp := ctx.Params("ref")

// 	if len(ap) > 0 {
// 		result := db.FirstOrCreate(u, &Miner{Address: ap})
// 		if result.RowsAffected == 1 {
// 			mon.Miners = append(mon.Miners, u)
// 		}
// 	}

// 	val := "%d%s__0"

// 	if len(rp) > 0 && u.ID != 0 {
// 		db.First(r, &Miner{Address: rp})
// 		if r.ID != 0 {
// 			u.ReferralID = r.ID
// 			err := db.Save(u).Error
// 			for err != nil {
// 				time.Sleep(time.Millisecond * 500)
// 				err = db.Save(u).Error
// 				log.Println(err)
// 			}
// 			val += "__" + r.Address
// 		}
// 	}

// 	dataTransaction(ap, &val, nil, nil)

// 	mr := &MineResponse{Success: true}
// 	ctx.JSON(200, mr)
// }

func minerView(ctx *macaron.Context) {
	height := getHeight()
	mr := &MinerResponse{}
	u := &Miner{}
	var miners []*Miner
	ap := ctx.Params("addr")

	if len(ap) > 0 {
		db.First(u, &Miner{Address: ap})
		if u.ID != 0 {
			mr.ID = u.ID
			mr.Exists = true
			mr.Address = u.Address
			mr.MiningHeight = u.MiningHeight
			mr.Height = height
			mr.MinedMobile = u.MinedMobile
			mr.MinedTelegram = u.MinedTelegram
			mr.TelegramId = u.TelegramId
			mr.AlphaSent = getAlphaSent(mr.Address)
			mr.Cycles = u.Cycles

			if u.TelegramId != 0 {
				mr.HasTelegram = true
			}

			db.Where("referral_id = ?", u.ID).Find(&miners).Count(&mr.Referred)

			db.Where("referral_id = ? AND mining_height > ?", u.ID, height-2880).Find(&miners).Count(&mr.Active)
		}
	}

	mr.Price = pc.AnotePrice

	ctx.JSON(200, mr)
}

func telegramMinerView(ctx *macaron.Context) {
	height := getHeight()
	mr := &MinerResponse{}
	u := &Miner{}
	var miners []*Miner
	tid := ctx.Params("tid")
	tidi, err := strconv.Atoi(tid)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	u = getMinerTel(int64(tidi))

	if u.ID != 0 {
		mr.ID = u.ID
		mr.Exists = true
		mr.Address = u.Address
		mr.MiningHeight = u.MiningHeight
		mr.Height = height
		mr.MinedMobile = u.MinedMobile
		mr.MinedTelegram = u.MinedTelegram
		mr.TelegramId = u.TelegramId

		if u.TelegramId != 0 {
			mr.HasTelegram = true
		}

		db.Where("referral_id = ?", u.ID).Find(&miners).Count(&mr.Referred)

		db.Where("referral_id = ? AND mining_height > ?", u.ID, height-2880).Find(&miners).Count(&mr.Active)
	}

	ctx.JSON(200, mr)
}

func withdrawView(ctx *macaron.Context) {
	mr := &MineResponse{}
	if strings.Contains(ctx.Req.RemoteAddr, "127.0.0.1") {
		u := &Miner{}
		tid := ctx.Params("tid")
		tidi, err := strconv.Atoi(tid)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		u = getMinerTel(int64(tidi))

		sendAssetTelegram(u.MinedTelegram-Fee, "", u.Address)

		u.MinedTelegram = 0
		db.Save(u)

		mon.loadMiners()
	} else {
		mr.Error = 1
		mr.Success = false
	}

	ctx.JSON(200, mr)
}

func saveTelegram(ctx *macaron.Context) {
	// var result *gorm.DB
	mr := &MineResponse{Success: true}
	ap := ctx.Params("addr")
	tids := ctx.Params("telegramid")
	tid, err := strconv.Atoi(tids)
	if err != nil {
		mr.Success = false
		mr.Error = 1
	}

	if strings.Contains(ctx.Req.RemoteAddr, "127.0.0.1") {
		m := &Miner{}
		db.First(m, &Miner{TelegramId: int64(tid)})

		if m.ID == 0 {
			db.FirstOrCreate(m, &Miner{TelegramId: int64(tid), Address: tids})
		}

		if strings.HasPrefix(ap, "3A") {
			m.Address = ap
		} else {
			refid, err := strconv.Atoi(ap)
			if err == nil && m.ReferralID == 0 && m.ID != uint(refid) {
				// result = db.FirstOrCreate(m, &Miner{TelegramId: int64(tid), Address: tids, ReferralID: uint(refid)})
				m.ReferralID = uint(refid)
			}
		}

		err := db.Save(m).Error
		if err != nil {
			// time.Sleep(time.Millisecond * 500)
			// err = db.Save(m).Error
			log.Println(err)
			// logTelegram(err.Error())
			mr.Success = false
			mr.Error = 3
		}
		// mon.loadMiners()
	} else {
		mr.Success = false
		mr.Error = 4
	}

	ctx.JSON(200, mr)
}

// func inviteView(ctx *macaron.Context) {
// 	var referred []*Miner
// 	mr := &MineResponse{Success: true}
// 	ap := ctx.Params("addr")
// 	height := getHeight()

// 	m := &Miner{}
// 	db.First(m, &Miner{Address: ap})

// 	db.Where("referral_id = ? AND mining_height < ?", m.ID, height-1440).Find(&referred)
// 	if time.Since(m.LastInvite) > (time.Hour * 24) {
// 		for _, r := range referred {
// 			if r.TelegramId != 0 {
// 				sendInvite(r)
// 			}
// 		}
// 		m.LastInvite = time.Now()
// 		err := db.Save(m).Error
// 		if err != nil {
// 			log.Println(err)
// 			logTelegram(err.Error())
// 		}
// 	} else {
// 		mr.Success = false
// 		mr.Error = 1
// 	}

// 	ctx.JSON(200, mr)
// }

func telegramMineView(ctx *macaron.Context) {
	ip := GetRealIP(ctx.Req.Request)
	h := getHeight()
	mr := &MineResponse{
		Success: true,
		Error:   0,
	}

	t := ctx.Params("tid")
	c := ctx.Params("code")

	if strings.Contains(ip, "127.0.0.1") {
		code := strings.TrimSpace(c)
		code = regexp.MustCompile(`[^0-9]+`).ReplaceAllString(code, "")

		log.Println(code)

		codeInt, err := strconv.Atoi(code)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
			mr.Success = false
			mr.Error = 2
		} else {
			tid, err := strconv.Atoi(t)
			log.Println(tid)
			if err != nil {
				log.Println(err)
				logTelegram(err.Error())
				mr.Success = false
				mr.Error = 2
			} else {
				if int(codeInt) == getMiningCode() {
					m := getMinerTel(int64(tid))
					if int64(h)-m.MiningHeight > 1409 {
						if m.MiningHeight > 0 {
							sendMined(m.Address, int64(h)-int64(m.MiningHeight))
							sendMinedTelegram(m.Address, int64(h)-int64(m.MiningHeight))
							m.Cycles++
							m.MiningTime = time.Now()
							m.MiningHeight = int64(h)
							m.BatteryNotification = true
							err := db.Save(m).Error
							if err != nil {
								log.Println(err)
								logTelegram(err.Error())
							}
							m.saveInBlockchain()
						} else {
							if strings.HasPrefix(m.Address, "3A") {
								sendAsset(Fee, "", m.Address)
							}
							m.MinedTelegram = Fee
							m.MiningTime = time.Now()
							m.Cycles = 1
							m.MiningHeight = int64(h)
							m.UpdatedApp = true
							m.BatteryNotification = true
							if m.Address == "" {
								m.Address = strconv.Itoa(int(m.TelegramId))
							}
							err := db.Save(m).Error
							if err != nil {
								log.Println(err)
								logTelegram(err.Error())
							}
							m.saveInBlockchain()
							// sendNotificationFirst(m)
						}
						// mon.loadMiners()
					} else {
						mr.Success = false
						mr.Error = 4
					}
				} else {
					mr.Success = false
					mr.Error = 3
				}
			}
		}
	} else {
		mr.Success = false
		mr.Error = 1
	}

	ctx.JSON(200, mr)
}
