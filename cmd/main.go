package main

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/fvbock/endless"

	apiHTTP "gin-items/api/http"
	"gin-items/library/setting"
)

func main() {
	router := apiHTTP.Init()

	endless.DefaultReadTimeOut = time.Duration(setting.Config().Server.ReadTimeout) * time.Second
	endless.DefaultWriteTimeOut = time.Duration(setting.Config().Server.WriteTimeout) * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.Config().Server.HttpPort)

	server := endless.NewServer(endPoint, router)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
