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

	m.Get("/mine/:address/:captchaid/:captcha/:code", mineView)
	m.Get("/mine/:address/:captchaid/:captcha/:code/:ref", mineView)
	m.Get("/new-captcha/:addr", newCaptchaView)

	return m
}
