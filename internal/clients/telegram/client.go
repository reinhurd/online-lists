package telegram

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"online-lists/internal/clients/yandex"
	"online-lists/internal/helpers"
)

func StartBot(tgToken string, yacl *yandex.Client) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var resp string
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			resp = "Hello, " + update.Message.From.UserName + "!" + " You said: " + update.Message.Text

			//TODO move to handler
			if update.Message.Text == "/headers" {
				res := helpers.GetCSVHeaders("internal/repository/SPISKEN.csv")
				resp = strings.Join(res, ", ")
			}
			if strings.Contains(update.Message.Text, "/add") {
				splStr := strings.Split(update.Message.Text, " ")
				err = helpers.InsertNewValueUnderHeader("internal/repository/SPISKEN.csv", splStr[1], splStr[2])
				if err != nil {
					fmt.Println(err)
				}
				resp = fmt.Sprintf("Added %s under %s", splStr[2], splStr[1])
			}
			if strings.Contains(update.Message.Text, "/ya_file") {
				yacl.GetYDFileByPath(os.Getenv("YDFILE"))
				helpers.ConvertToCSV()
				resp = "File downloaded and converted to CSV"
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
