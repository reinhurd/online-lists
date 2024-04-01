package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"online-lists/internal/clients/telegram"
	"online-lists/internal/clients/yandex"
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
	//start telegram bot

	go func() {
		telegram.StartBot(os.Getenv("TG_SECRET_KEY"), yaClient)
	}()

	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
