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
		log.Println(err)
		logTelegram(err.Error())
		savedHeight = 0
	} else {
		sh := parseItem(minerData.(string), 1)
		savedRef = parseItem(minerData.(string), 3)
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

	if pr.Error == 0 && (height-savedHeight > 1410) && !sendTelegramNotification(addr, height, savedHeight) {
		pr.Success = false
		pr.Error = 3
	}

	if pr.Error == 0 && (height-savedHeight > 1410) {
		log.Println(fmt.Sprintf("%s %s", addr, ip))

		encIp := EncryptMessage(ip)

		newMinerData := updateItem(minerData.(string), height, 1)
		newMinerData = updateItem(newMinerData, encIp, 2)

		if savedRef != nil && len(savedRef.(string)) > 0 {
			newMinerData = updateItem(newMinerData, savedRef.(string), 3)
		} else if len(ref) > 0 {
			newMinerData = updateItem(newMinerData, ref, 3)
		}

		dataTransaction(addr, &newMinerData, nil, nil)

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

type ImageResponse struct {
	Image string `json:"image"`
	Id    string `json:"id"`
}
