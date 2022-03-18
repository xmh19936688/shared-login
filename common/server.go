package common

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RunServer(gs *gin.Engine, port string) *http.Server {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: gs,
	}

	go func(srv *http.Server) {
		log.Println("listening:", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln("failed to serve:" + err.Error())
		}
	}(srv)

	return srv
}

func StopServer(srv *http.Server) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	srv.Shutdown(ctx)
}
