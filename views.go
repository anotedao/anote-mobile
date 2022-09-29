package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-macaron/captcha"
	macaron "gopkg.in/macaron.v1"
)

func mineView(ctx *macaron.Context, cpt *captcha.Captcha) {
	pr := &MineResponse{
		Success: true,
		Error:   0,
	}

	addr := ctx.Params("address")
	cpid := ctx.Params("captchaid")
	cp := ctx.Params("captcha")
	code := ctx.Params("code")

	codeInt, err := strconv.Atoi(code)
	if err != nil {
		log.Println(err)
	}

	if !cpt.Verify(cpid, cp) {
		pr.Success = false
		pr.Error = 1
	}

	if int(codeInt) != getMiningCode() {
		pr.Success = false
		pr.Error = 2
	}

	var savedHeight int64
	sh, err := getData(addr, nil)
	if err != nil {
		log.Println(err)
		savedHeight = 0
	} else {
		savedHeight = sh.(int64)
	}
	height := int64(getHeight())

	if pr.Error == 0 && (height-savedHeight > 1440) && !sendTelegramNotification(addr) {
		pr.Success = false
		pr.Error = 3
	}

	if pr.Error == 0 && (height-savedHeight > 1440) {
		dataTransaction(addr, nil, &height, nil)
		if sh != nil && height-savedHeight <= 2880 {
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
