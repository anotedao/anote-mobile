package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-macaron/captcha"
	macaron "gopkg.in/macaron.v1"
)

func mineView(ctx *macaron.Context, cpt *captcha.Captcha) {
	var savedHeight int64
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

	log.Println(ref)

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
		if sh != nil {
			savedHeight = sh.(int64)
			savedHeight = 0
		}
	}

	log.Println(savedHeight)

	if pr.Error == 0 && (height-savedHeight > 1440) && !sendTelegramNotification(addr) {
		pr.Success = false
		pr.Error = 3
	}

	if pr.Error == 0 && (height-savedHeight > 1440) {
		newMinerData := updateItem(minerData.(string), height, 1)
		if len(ref) > 0 {
			newMinerData = updateItem(newMinerData, ref, 2)
		}
		dataTransaction(addr, &newMinerData, nil, nil)
		if height-savedHeight <= 2880 {
			go sendMined(addr)
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
