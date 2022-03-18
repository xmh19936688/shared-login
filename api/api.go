package api

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"xmh.shared-login/model"
	"xmh.shared-login/setting"
)

var cacheMap sync.Map

func checkAuth(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	log.Println("auth-header:", token)
	if token != setting.Config.Auth.Secret {
		c.JSON(http.StatusUnauthorized, "invalid token")
		c.Abort()
	}
}

func auth(c *gin.Context) {
	var form model.Identity
	if err := binding.JSON.Bind(c.Request, &form); err != nil {
		c.JSON(http.StatusBadRequest, "invalid input")
		return
	}

	if len(form.Password) == 0 && len(form.Token) == 0 {
		c.JSON(http.StatusUnauthorized, "invalid input")
		return
	}

	if len(form.Password) > 0 && form.Password != setting.Config.Auth.Secret {
		c.JSON(http.StatusUnauthorized, "invalid password")
		return
	}
	if len(form.Token) > 0 && form.Token != setting.Config.Auth.Secret {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	if len(form.Username) > 0 && len(form.Password) > 0 {
		c.JSON(http.StatusOK, model.Identity{Token: setting.Config.Auth.Secret})
	}
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", "")
}

func set(c *gin.Context) {
	var form model.KV
	if err := binding.JSON.Bind(c.Request, &form); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, "invalid input")
		return
	}

	cacheMap.Store(form.Key, form.Value)
	c.JSON(http.StatusOK, "success")
}

func get(c *gin.Context) {
	var form model.KV
	if err := binding.JSON.Bind(c.Request, &form); err != nil {
		c.JSON(http.StatusBadRequest, "invalid input")
		return
	}

	value, ok := cacheMap.Load(form.Key)
	if !ok {
		c.JSON(http.StatusBadRequest, "not found")
		return
	}

	c.JSON(http.StatusOK, value)
}

func Init(gs *gin.Engine) {
	gs.GET("/", index)
	gs.POST("/auth", auth)
	gs.POST("/set", checkAuth, set)
	gs.POST("/get", checkAuth, get)
}
