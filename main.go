package main

import (
	"fmt"
	"net/http"
	"time"

	apiHTTP "gin-items/api/http"
	"gin-items/library/setting"
)

func main()  {

	router := apiHTTP.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.Config().Server.HttpPort),
		Handler:        router,
		ReadTimeout:    time.Duration(setting.Config().Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(setting.Config().Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
