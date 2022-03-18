package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"

	"xmh.shared-login/api"
	"xmh.shared-login/common"
)

func main() {
	fmt.Println("server")
	gin.SetMode(gin.DebugMode)
	gs := gin.Default()
	gs.LoadHTMLGlob("static/*")
	gs.Static("/static", "static")

	api.Init(gs)
	srv := common.RunServer(gs, "9000")

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(interrupt, os.Kill)
	<-interrupt

	common.StopServer(srv)
}
