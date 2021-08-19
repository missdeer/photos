package main

import (
	"os"
	_ "github.com/missdeer/photos/routers"

	"github.com/beego/beego"
)

const (
	hostVar = "APP_HOST"
	portVar = "APP_PORT"
)

func main() {
	var port string
	if port = os.Getenv(portVar); port == "" {
		port = "8899"
	}
	beego.Run(beego.BConfig.Listen.HTTPAddr + ":" + port)
}
