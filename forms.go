package main

type MineForm struct {
	Code      string `form:"code" binding:"Required"`
	Captcha   string `form:"captcha" binding:"Required"`
	CaptchaId string `form:"captcha_id" binding:"Required"`
}
