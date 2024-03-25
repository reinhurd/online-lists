package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"online-lists/internal/clients/telegram"
	"online-lists/internal/clients/yandex"
	"online-lists/internal/helpers"
)

var db = make(map[string]string)

func setupRouter(yaCl *yandex.Client) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	r.GET("/ya_list", func(c *gin.Context) {
		yaCl.GetYDList()
	})
	r.GET("/ya_file", func(c *gin.Context) {
		yaCl.GetYDFileByPath(os.Getenv("YDFILE"))
		helpers.ConvertToCSV()
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
	r := setupRouter(yaClient)
	//start telegram bot
	telegram.StartBot(os.Getenv("TG_SECRET_KEY"))

	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
