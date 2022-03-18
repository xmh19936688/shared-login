package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"xmh.shared-login/setting"

	"github.com/gin-gonic/gin"

	"xmh.shared-login/common"
	"xmh.shared-login/model"
	"xmh.shared-login/utils"
)

var token = ""

// 浏览器带有tk，且本地tk为空，则更新本地tk
// 浏览器不带tk，且本地tk非空，则返回本地tk
func sync(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Header("Access-Control-Expose-Headers", "Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")

	tk := c.Request.Header.Get("Authorization")
	log.Println("auth-header:", tk)
	// 如果浏览器有token，更新为浏览器token
	if len(tk) > 0 {
		token = tk
		fmt.Println("token updated:", token)
		return
	}

	// 如果本地有token，返回本地token
	if len(token) > 0 {
		c.JSON(http.StatusOK, model.Identity{Token: setting.Config.Auth.Secret})
	}
}

func initApi(gs *gin.Engine) {
	gs.POST("/sync", sync)
	gs.OPTIONS("/sync", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
	})
}

func handle() {
	for {
		fmt.Println("input cmd (login, token, exit):")
		var str string
		_, err := fmt.Scanln(&str)
		if err != nil {
			fmt.Println(err.Error())
		}
		switch str {
		case "login":
			login()
		case "token":
			fmt.Println("token is:", token)
		case "exit":
			os.Exit(0)
		}
	}
}

func login() {
	urlStr := "http://localhost:9000/auth"
	data := model.Identity{
		Username: "admin",
		Password: "admin",
	}
	bs, err := json.Marshal(data)
	if err != nil {
		utils.EP(err)
		return
	}

	request, err := http.NewRequest("POST", urlStr, bytes.NewReader(bs))
	if err != nil {
		utils.EP(err)
		return
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		utils.EP(err)
		return
	}

	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.EP(err)
		return
	}

	err = json.Unmarshal(bs, &data)
	if err != nil {
		utils.EP(err)
		return
	}
	token = data.Token
	fmt.Println("login succeed. token:" + token)
}

func main() {
	fmt.Println("client")
	gin.SetMode(gin.DebugMode)
	gs := gin.Default()

	initApi(gs)
	srv := common.RunServer(gs, "9999")

	go handle()

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(interrupt, os.Kill)
	<-interrupt

	common.StopServer(srv)
}
