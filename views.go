package main

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-macaron/captcha"
	macaron "gopkg.in/macaron.v1"
)

func mineView(ctx *macaron.Context, cpt *captcha.Captcha) {
	height := int64(getHeight())

	pr := &MineResponse{
		Success: true,
		Error:   0,
	}

	addr := ctx.Params("address")
	cpid := ctx.Params("captchaid")
	cp := ctx.Params("captcha")
	code := ctx.Params("code")
	ip := GetRealIP(ctx.Req.Request)

	miner := getMinerOrCreate(addr)
	savedHeight := miner.MiningHeight

	code = strings.TrimSpace(code)
	code = regexp.MustCompile(`[^0-9]+`).ReplaceAllString(code, "")

	codeInt, err := strconv.Atoi(code)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		pr.Success = false
		pr.Error = 2
	}

	if !cpt.Verify(cpid, cp) {
		pr.Success = false
		pr.Error = 1
	}

	if int(codeInt) != getMiningCode() {
		pr.Success = false
		pr.Error = 2
	}

	if pr.Error == 0 && countIP(ip) > 3 {
		pr.Success = false
		pr.Error = 4
	}

	if !strings.HasPrefix(addr, "3A") {
		pr.Success = false
		pr.Error = 5
	}

	if pr.Error == 0 && (height-miner.MiningHeight > 1410) {
		log.Println(fmt.Sprintf("%s %s", addr, ip))

		miner.clearIps()
		miner.saveIp(ip)

		if savedHeight > 0 {
			go sendMined(addr, height-savedHeight)
			go func() {
				time.Sleep(time.Second * 30)
				checkConfirmation(addr)
			}()
		} else {
			go sendMinedFirst(addr)
			miner.PingCount = 1
			miner.MiningTime = time.Now()
			miner.MiningHeight = height
			miner.UpdatedApp = true
			miner.BatteryNotification = true
			err := db.Save(miner).Error
			for err != nil {
				time.Sleep(time.Millisecond * 500)
				err = db.Save(miner).Error
			}
			miner.saveInBlockchain()
			sendNotificationFirst(miner)
		}
	}

	ctx.Resp.Header().Add("Access-Control-Allow-Origin", "*")
	ctx.JSON(200, pr)
}

func newCaptchaView(ctx *macaron.Context, cpt *captcha.Captcha) {
	c, err := cpt.CreateCaptcha()
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ir := &ImageResponse{
		Id:    c,
		Image: fmt.Sprintf("%s/captcha/%s.png", conf.Host, c),
	}

	ctx.Resp.Header().Add("Access-Control-Allow-Origin", "*")
	ctx.JSON(200, ir)
}

type MineResponse struct {
	Success bool `json:"success"`
	Error   int  `json:"error"`
}

type MinerResponse struct {
	Address      string `json:"address"`
	Referred     int64  `json:"referred"`
	Active       int64  `json:"active"`
	Confirmed    int64  `json:"confirmed"`
	HasTelegram  bool   `json:"has_telegram"`
	MiningHeight int64  `json:"mining_height"`
	Height       uint64 `json:"height"`
	Exists       bool   `json:"exists"`
}

type MinePingResponse struct {
	Success       bool `json:"success"`
	CycleFinished bool `json:"cycle_finished"`
	Error         int  `json:"error"`
	Health        int  `json:"health"`
}

type HealthResponse struct {
	Health     int  `json:"health"`
	UpdatedApp bool `json:"updated_app"`
}

type ImageResponse struct {
	Image string `json:"image"`
	Id    string `json:"id"`
}

func minePingView(ctx *macaron.Context) {
	a := ctx.Params("address")
	apk := ctx.Params("apk")
	ip := GetRealIP(ctx.Req.Request)

	mr := &MinePingResponse{Success: true}
	mr.CycleFinished = false

	height := int64(mon.Height)

	miner := getMiner(a)

	if miner.ID == 0 {
		mr.Success = false
		mr.Error = 1
		mr.CycleFinished = true
	} else if !strings.HasPrefix(a, "3A") {
		mr.Success = false
		mr.Error = 2
		mr.CycleFinished = true
	} else {
		if height-miner.MiningHeight >= 1410 {
			mr.CycleFinished = true
		}

		if time.Since(miner.LastPing) > time.Second*55 {
			miner.saveIp(ip)
			minerPing(miner)

			if apk == conf.APK {
				miner.UpdatedApp = true
			} else {
				miner.UpdatedApp = false
			}

			db.Save(miner)
		}

		if mr.CycleFinished {
			if mon.isSending(miner, 1410) {
				sendNotificationEnd(miner)
				log.Printf("Notification: %s", miner.Address)
			}
		}
	}

	// m := time.Since(miner.MiningTime).Minutes()
	mr.Health = int((math.Floor(getIpFactor(miner, true, uint64(height), 0)*100) / 100) * 100)

	if mr.Health > 100 {
		mr.Health = 100
	} else if mr.Health < 0 {
		mr.Health = 0
	}

	log.Println("Ping: " + a + " " + ip + " Health: " + strconv.Itoa(mr.Health) + " IPs: " + strconv.Itoa(int(db.Model(miner).Association("IpAddresses").Count())) + " Pings: " + strconv.Itoa(int(miner.PingCount)))

	ctx.JSON(200, mr)
}

func healthView(ctx *macaron.Context) {
	a := ctx.Params("address")
	height := getHeight()

	hr := &HealthResponse{}

	miner := getMiner(a)

	hr.Health = int((math.Floor(getIpFactor(miner, true, uint64(height), 2)*100) / 100) * 100)

	if hr.Health > 100 {
		hr.Health = 100
	} else if hr.Health < 0 {
		hr.Health = 0
	}

	hr.UpdatedApp = miner.UpdatedApp

	ctx.JSON(200, hr)
}

func statsView(ctx *macaron.Context) {
	sr := getStats()
	ctx.JSON(200, sr)
}

func minerPing(miner *Miner) {
	miner.PingCount++
	miner.LastPing = time.Now()
	db.Save(miner)
}

func newUserView(ctx *macaron.Context) {
	u := &Miner{}
	r := &Miner{}

	ap := ctx.Params("addr")
	rp := ctx.Params("ref")

	if len(ap) > 0 {
		result := db.FirstOrCreate(u, &Miner{Address: ap})
		if result.RowsAffected == 1 {
			mon.Miners = append(mon.Miners, u)
		}
	}

	val := "%d%s__0"

	if len(rp) > 0 && u.ID != 0 {
		db.First(r, &Miner{Address: rp})
		if r.ID != 0 {
			u.ReferralID = r.ID
			db.Save(u)
			val += "__" + r.Address
		}
	}

	dataTransaction(ap, &val, nil, nil)

	mr := &MineResponse{Success: true}
	ctx.JSON(200, mr)
}

func minerView(ctx *macaron.Context) {
	height := getHeight()
	mr := &MinerResponse{}
	u := &Miner{}
	var miners []*Miner
	ap := ctx.Params("addr")

	if len(ap) > 0 {
		db.First(u, &Miner{Address: ap})
		if u.ID != 0 {
			mr.Exists = true
			mr.Address = u.Address
			mr.MiningHeight = u.MiningHeight
			mr.Height = height

			if u.TelegramId != 0 {
				mr.HasTelegram = true
			}

			db.Where("referral_id = ?", u.ID).Find(&miners).Count(&mr.Referred)

			db.Where("referral_id = ? AND mining_height > ?", u.ID, height-2880).Find(&miners).Count(&mr.Active)

			db.Where("referral_id = ? AND mining_height > ? AND confirmed = true", u.ID, height-2880).Find(&miners).Count(&mr.Confirmed)
		}
	}

	ctx.JSON(200, mr)
}

func saveTelegram(ctx *macaron.Context) {
	mr := &MineResponse{Success: true}
	ap := ctx.Params("addr")
	tids := ctx.Params("telegramid")
	tid, err := strconv.Atoi(tids)
	if err != nil {
		mr.Success = false
		mr.Error = 1
	}

	m := &Miner{}
	result := db.FirstOrCreate(m, &Miner{Address: ap})
	if result.RowsAffected == 1 {
		mon.Miners = append(mon.Miners, m)
	}

	if strings.Contains(ctx.Req.RemoteAddr, "127.0.0.1") {
		// if m.ID == 0 {
		// 	mr.Success = false
		// 	mr.Error = 2
		// } else if m.TelegramId != 0 {
		// 	mr.Success = false
		// 	mr.Error = 3
		// } else {
		// 	m.TelegramId = int64(tid)
		// 	db.Save(m)
		// }
		m.TelegramId = int64(tid)
		db.Save(m)
	} else {
		mr.Success = false
		mr.Error = 4
	}

	ctx.JSON(200, mr)
}

func inviteView(ctx *macaron.Context) {
	var referred []*Miner
	mr := &MineResponse{Success: true}
	ap := ctx.Params("addr")
	height := getHeight()

	m := &Miner{}
	db.First(m, &Miner{Address: ap})

	db.Where("referral_id = ? AND mining_height < ?", m.ID, height-1440).Find(&referred)
	if time.Since(m.LastInvite) > (time.Hour * 24) {
		go func() {
			for _, r := range referred {
				if r.TelegramId != 0 {
					sendInvite(r)
				}
			}
		}()
		m.LastInvite = time.Now()
		db.Save(m)
	} else {
		mr.Success = false
		mr.Error = 1
	}

	ctx.JSON(200, mr)
}

func telegramMineView(ctx *macaron.Context) {
	ip := GetRealIP(ctx.Req.Request)
	mr := &MineResponse{
		Success: true,
		Error:   0,
	}

	t := ctx.Params("tid")
	c := ctx.Params("code")

	log.Println(t)
	log.Println(c)

	if strings.Contains(ip, "127.0.0.1") {

	} else {
		mr.Success = false
		mr.Error = 1
	}

	ctx.JSON(200, mr)
}
