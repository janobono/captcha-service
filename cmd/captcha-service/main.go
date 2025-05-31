package main

import (
	"github.com/janobono/captcha-service/internal/config"
	"github.com/janobono/captcha-service/internal/server"
)

func main() {
	serverConfig := config.InitConfig()
	app := server.NewServer(serverConfig)
	app.Start()
}
