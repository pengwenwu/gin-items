package main

import (
	"fmt"
	"net/http"

	apiHTTP "gin-items/api/http"
	"gin-items/library/setting"
)

func main()  {

	router := apiHTTP.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
