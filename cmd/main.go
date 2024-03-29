package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"online-lists/internal/clients/telegram"
	"online-lists/internal/clients/yandex"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	if err := godotenv.Load("secret.env"); err != nil {
		panic(err)
	}
	YA_ID := os.Getenv("YANDEX_TOKEN")
	restyCl := resty.New()
	yaClient := yandex.NewClient(restyCl, YA_ID)

	r := setupRouter()
	//start telegram bot
	telegram.StartBot(os.Getenv("TG_SECRET_KEY"), yaClient)

	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
