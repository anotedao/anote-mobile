package main

import (
	"fmt"
	"log"

	"github.com/go-macaron/captcha"
	macaron "gopkg.in/macaron.v1"
)

func mineView(ctx *macaron.Context) {
	pr := &PingResponse{
		Success: true,
	}

	ctx.JSON(200, pr)
}

func newCaptchaView(ctx *macaron.Context, cpt *captcha.Captcha) {
	c, err := cpt.CreateCaptcha()
	if err != nil {
		log.Println(err)
	}

	ir := &ImageResponse{
		Id:    c,
		Image: fmt.Sprintf("https://mobile.anote.digital/captcha/%s.png", c),
	}

	ctx.Resp.Header().Add("Access-Control-Allow-Origin", "*")
	ctx.JSON(200, ir)
}

type PingResponse struct {
	Success bool `json:"success"`
}

type ImageResponse struct {
	Image string `json:"image"`
	Id    string `json:"id"`
}
