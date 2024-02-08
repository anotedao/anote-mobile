package main

import (
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	mac := macaron.Classic()

	mac.Use(macaron.Renderer())
	mac.Use(cache.Cacher())
	mac.Use(captcha.Captchaer(
		captcha.Options{
			ChallengeNums: 5,
		},
	))

	mac.Get("/save-telegram/:addr/:telegramid", saveTelegram)
	// mac.Get("/mine/:address/:captchaid/:captcha/:code", mineView)
	// mac.Get("/mine/:address/:captchaid/:captcha/:code/:ref", mineView)
	// mac.Get("/new-captcha/:addr", newCaptchaView)
	// mac.Get("/new-user/:addr", newUserView)
	// mac.Get("/new-user/:addr/:ref", newUserView)
	mac.Get("/miner/:addr", minerView)
	mac.Get("/tminer/:tid", telegramMinerView)
	mac.Get("/withdraw/:tid", withdrawView)
	// mac.Get("/health/:address", healthView)
	mac.Get("/stats", statsView)
	// mac.Get("/invite/:addr", inviteView)
	mac.Get("/telegram-mine/:tid/:code", telegramMineView)

	return mac
}
