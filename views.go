package main

import (
	macaron "gopkg.in/macaron.v1"
)

func mineView(ctx *macaron.Context) {
	pr := &PingResponse{
		Success: true,
	}

	ctx.JSON(200, pr)
}

type PingResponse struct {
	Success bool `json:"success"`
}
