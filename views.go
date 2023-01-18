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

	if pr.Error == 0 && (height-miner.MiningHeight > 1410) {
		log.Println(fmt.Sprintf("%s %s", addr, ip))

		miner.IP = ip
		miner.IP2 = ""
		miner.IP3 = ""
		miner.IP4 = ""
		miner.IP5 = ""

		if savedHeight > 0 {
			go sendMined(addr, height-savedHeight)
			go func() {
				time.Sleep(time.Second * 30)
				checkConfirmation(addr)
			}()
		} else {
			miner.PingCount = 1
			miner.MiningTime = time.Now()
			miner.MiningHeight = height
		}

		db.Save(miner)

		miner.saveInBlockchain()
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

type MinePingResponse struct {
	Success       bool `json:"success"`
	CycleFinished bool `json:"cycle_finished"`
	Error         int  `json:"error"`
	Health        int  `json:"health"`
}

type HealthResponse struct {
	Health int `json:"health"`
}

type ImageResponse struct {
	Image string `json:"image"`
	Id    string `json:"id"`
}

func minePingView(ctx *macaron.Context) {
	a := ctx.Params("address")
	ip := GetRealIP(ctx.Req.Request)

	mr := &MinePingResponse{Success: true}
	mr.CycleFinished = false

	height := int64(getHeight())

	miner := getMiner(a)

	if miner.ID == 0 {
		mr.Success = false
		mr.Error = 1
	} else {
		if height-miner.MiningHeight > 1410 {
			mr.CycleFinished = true
		}

		if time.Since(miner.LastPing) > time.Second*59 {
			if ip == miner.IP {
				minerPing(miner)
			} else if len(miner.IP2) == 0 || miner.IP2 == ip {
				miner.IP2 = ip
				minerPing(miner)
			} else if len(miner.IP3) == 0 || miner.IP3 == ip {
				miner.IP3 = ip
				minerPing(miner)
			} else if len(miner.IP4) == 0 || miner.IP4 == ip {
				miner.IP4 = ip
				minerPing(miner)
			} else if len(miner.IP5) == 0 || miner.IP5 == ip {
				miner.IP5 = ip
				minerPing(miner)
			}
		}
	}

	// m := time.Since(miner.MiningTime).Minutes()
	mr.Health = int(getIpFactor(miner) * 100)

	if mr.Health > 100 {
		mr.Health = 100
	} else if mr.Health < 0 {
		mr.Health = 0
	}

	log.Println("Ping: " + a + " " + ip + " " + strconv.Itoa(int(getIpFactor(miner)*100)))

	ctx.JSON(200, mr)
}

func healthView(ctx *macaron.Context) {
	a := ctx.Params("address")

	hr := &HealthResponse{}

	miner := getMiner(a)

	hr.Health = int(getIpFactor(miner) * 100)

	if hr.Health > 100 {
		hr.Health = 100
	} else if hr.Health < 0 {
		hr.Health = 0
	}

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
