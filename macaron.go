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
			Width:         160,
			Height:        60,
		}))

	m.Get("/mine/:addr", mineView)
	m.Get("/new-captcha/:addr", newCaptchaView)

	return m
}
