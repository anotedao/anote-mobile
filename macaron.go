package main

import (
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	m := macaron.Classic()

	m.Use(macaron.Renderer())
	m.Use(cache.Cacher())
	m.Use(captcha.Captchaer(
		captcha.Options{
			ChallengeNums: 5,
		},
	))

	m.Get("/save-telegram/:addr/:telegramid", saveTelegram)
	m.Get("/mine/:address/:captchaid/:captcha/:code", mineView)
	m.Get("/mine/:address/:captchaid/:captcha/:code/:ref", mineView)
	m.Get("/new-captcha/:addr", newCaptchaView)
	m.Get("/new-user/:addr", newUserView)
	m.Get("/new-user/:addr/:ref", newUserView)
	m.Get("/mine/:address", minePingView)
	m.Get("/miner/:addr", minerView)
	m.Get("/mine/:address/:apk", minePingView)
	m.Get("/health/:address", healthView)
	m.Get("/stats", statsView)
	m.Get("/invite/:addr", inviteView)

	return m
}
