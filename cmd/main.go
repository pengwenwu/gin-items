package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiHTTP "gin-items/api/http"
	"gin-items/library/setting"
)

func main() {
	router := apiHTTP.Init()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.Config().Server.HttpPort),
		Handler:        router,
		ReadTimeout:    time.Duration(setting.Config().Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(setting.Config().Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// 服务连接
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
