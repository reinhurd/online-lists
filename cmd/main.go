package main

import (
	"context"
	"errors"
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
	//todo create range handler
	r.GET("/ya_list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.GetYaList(),
		})
	})
	r.GET("/headers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.GetHeaders(),
		})
	})
	r.GET("/set_csv", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.SetDefaultCsv(c.Query("filename")),
		})
	})
	r.GET("/list_csv", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.ListCsv(),
		})
	})
	r.GET("/add", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.Add(c.Query("header"), c.Query("value")),
		})
	})
	r.GET("/ya_file", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.YAFile(c.Query("filename")),
		})
	})
	r.GET("/ya_upload", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.YAUpload(c.Query("filename")),
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
	tgbot, err := telegram.StartBot(os.Getenv("TG_SECRET_KEY"), svc, true)
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
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		tgbot.SendToLastChat("Service is shutting down with error")
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
		tgbot.SendToLastChat("Service is shutting down by timeout")
	}
	log.Println("Server exiting")
}
