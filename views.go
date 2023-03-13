package main

import (
	"fmt"
	"log"
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
			db.Save(miner)
			miner.saveInBlockchain()
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
	Address     string `json:"address"`
	Referred    int64  `json:"referred"`
	Active      int64  `json:"active"`
	HasTelegram bool   `json:"has_telegram"`
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

	height := int64(getHeight())

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
		if height-miner.MiningHeight > 1410 {
			mr.CycleFinished = true
		}

		if time.Since(miner.LastPing) > time.Second*55 {
			// if ip == miner.IP {
			// 	minerPing(miner)
			// } else if len(miner.IP2) == 0 || miner.IP2 == ip {
			// 	miner.IP2 = ip
			// 	minerPing(miner)
			// } else if len(miner.IP3) == 0 || miner.IP3 == ip {
			// 	miner.IP3 = ip
			// 	minerPing(miner)
			// } else if len(miner.IP4) == 0 || miner.IP4 == ip {
			// 	miner.IP4 = ip
			// 	minerPing(miner)
			// } else if len(miner.IP5) == 0 || miner.IP5 == ip {
			// 	miner.IP5 = ip
			// 	minerPing(miner)
			// }
			miner.saveIp(ip)
			minerPing(miner)

			if apk == conf.APK {
				miner.UpdatedApp = true
			} else {
				miner.UpdatedApp = false
			}

			db.Save(miner)
		}
	}

	// m := time.Since(miner.MiningTime).Minutes()
	mr.Health = int(getIpFactor(miner, true, uint64(height)) * 100)

	if mr.Health > 100 {
		mr.Health = 100
	} else if mr.Health < 0 {
		mr.Health = 0
	}

	log.Println("Ping: " + a + " " + ip + " Health: " + strconv.Itoa(int(getIpFactor(miner, true, uint64(height))*100)) + " IPs: " + strconv.Itoa(int(db.Model(miner).Association("IpAddresses").Count())) + " Pings: " + strconv.Itoa(int(miner.PingCount)))

	ctx.JSON(200, mr)
}

func healthView(ctx *macaron.Context) {
	a := ctx.Params("address")
	height := getHeight()

	hr := &HealthResponse{}

	miner := getMiner(a)

	hr.Health = int(getIpFactor(miner, true, height) * 100)

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
		db.FirstOrCreate(u, &Miner{Address: ap})
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
		mr.Address = u.Address
	}

	if u.TelegramId != 0 {
		mr.HasTelegram = true
	}

	db.Where("referral_id = ?", u.ID).Find(&miners).Count(&mr.Referred)

	db.Where("referral_id = ? AND mining_height > ?", u.ID, height-2880).Find(&miners).Count(&mr.Active)

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
	db.FirstOrCreate(m, &Miner{Address: ap})

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
