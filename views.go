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
	var savedHeight int64
	var savedRef interface{}
	height := int64(getHeight())
	pr := &MineResponse{
		Success: true,
		Error:   0,
	}

	addr := ctx.Params("address")
	cpid := ctx.Params("captchaid")
	cp := ctx.Params("captcha")
	code := ctx.Params("code")
	ref := ctx.Params("ref")
	ip := GetRealIP(ctx.Req.Request)

	miner := getMinerOrCreate(addr)

	log.Println(ref)

	code = strings.TrimSpace(code)
	code = regexp.MustCompile(`[^0-9]+`).ReplaceAllString(code, "")

	codeInt, err := strconv.Atoi(code)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if !cpt.Verify(cpid, cp) {
		pr.Success = false
		pr.Error = 1
	}

	if int(codeInt) != getMiningCode() {
		pr.Success = false
		pr.Error = 2
	}

	minerData, err := getData(addr, nil)
	if err != nil {
		savedHeight = 0
		md := "%d%s__0"
		minerData = md
		dataTransaction(addr, &md, nil, nil)
	} else {
		sh := parseItem(minerData.(string), 0)
		savedRef = parseItem(minerData.(string), 1)
		if sh != nil {
			savedHeight = int64(sh.(int))
		} else {
			savedHeight = 0
		}
	}

	log.Println(savedHeight)

	if pr.Error == 0 && countIP(ip) > 3 {
		pr.Success = false
		pr.Error = 4
	}

	if pr.Error == 0 && (height-savedHeight > 1410) {
		log.Println(fmt.Sprintf("%s %s", addr, ip))
		newMinerData := updateItem(minerData.(string), height, 0)

		if savedRef != nil && len(savedRef.(string)) > 0 {
			newMinerData = updateItem(newMinerData, savedRef.(string), 1)
		} else if len(ref) > 0 {
			newMinerData = updateItem(newMinerData, ref, 1)
		}

		log.Println(newMinerData)

		dataTransaction(addr, &newMinerData, nil, nil)
		miner.MiningHeight = height
		miner.IP = ip
		miner.IP2 = ""
		miner.IP3 = ""
		miner.MiningTime = time.Now()
		db.Save(miner)

		if savedHeight > 0 {
			go sendMined(addr, height-savedHeight)
			go func() {
				time.Sleep(time.Second * 30)
				checkConfirmation(addr)
			}()
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

type MinePingResponse struct {
	Success       bool `json:"success"`
	CycleFinished bool `json:"cycle_finished"`
	Error         int  `json:"error"`
	Health        int  `json:"health"`
}

type ImageResponse struct {
	Image string `json:"image"`
	Id    string `json:"id"`
}

func minePingView(ctx *macaron.Context) {
	a := ctx.Params("address")
	ip := GetRealIP(ctx.Req.Request)

	log.Println("Ping: " + a + " " + ip)

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

		if time.Since(miner.LastPing) > time.Second*55 {
			if ip == miner.IP {
				minerPing(miner)
			} else if len(miner.IP2) == 0 || miner.IP2 == ip {
				miner.IP2 = ip
				minerPing(miner)
			} else if len(miner.IP3) == 0 || miner.IP3 == ip {
				miner.IP3 = ip
				minerPing(miner)
			}
		}
	}

	s := time.Since(miner.MiningTime).Seconds()
	mr.Health = int(float64(miner.PingCount) / float64(int64(s)/60) * 100)

	if mr.Health > 100 {
		mr.Health = 100
	}

	ctx.JSON(200, mr)
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
