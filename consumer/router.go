package main

import (
	"PowderDev/kafka-top/websockets"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func getRouter() *router.Router {
	r := router.New()

	r.GET("/ws", func(c *fasthttp.RequestCtx) {
		websockets.ServeWs(c, wsHub)
	})

	return r
}
