package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"online-lists/internal/clients/telegram"
	"online-lists/internal/clients/yandex"
	"online-lists/internal/helpers"
	"online-lists/internal/service"
)

func setupRouter(svc *service.Service) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/ya_list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.GetYaList(),
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, helpers.GetHomeTemplate())
	})

	return r
}

func main() {
	if err := godotenv.Load("secret.env"); err != nil {
		panic(err)
	}
	YaId := os.Getenv("YANDEX_TOKEN")
	restyCl := resty.New()
	yaClient := yandex.NewClient(restyCl, YaId)
	svc := service.NewService(yaClient)

	r := setupRouter(svc)
	tgbot, err := telegram.StartBot(os.Getenv("TG_SECRET_KEY"), yaClient, true)
	if err != nil {
		panic(err)
	}

	go func() {
		updChan := tgbot.GetUpdatesChan()
		err = tgbot.HandleUpdate(updChan)
		if err != nil {
			panic(err)
		}
	}()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	shutdownTimeout := 15 * time.Second
	shutdown := make(chan os.Signal, 1)
	endShutdown := make(chan struct{})
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer func(os.Signal) {
		log.Println("received signal, start shutdown")
		//todo deal with infinite shutdown
		srv.Shutdown(ctx)
		tgbot.SendToLastChat("Service is shutting down")
		log.Println("received signal, end shutdown")
		endShutdown <- struct{}{}
	}(<-shutdown)

	select {
	case <-endShutdown:
		log.Println("shuthown end, goodbye")
		os.Exit(0)
	case <-time.After(shutdownTimeout):
		log.Println("shutdhown timeout, goodbye")
		tgbot.SendToLastChat("Service is shutting down")

		os.Exit(0)
	}
}
