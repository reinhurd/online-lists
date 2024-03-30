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

var defaultCsvName string

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
				res := helpers.GetCSVHeaders("internal/repository/" + defaultCsvName)
				resp = strings.Join(res, ", ")
			}
			if strings.Contains(update.Message.Text, "/set_csv") {
				splStr := strings.Split(update.Message.Text, " ")
				defaultCsvName = splStr[1]
				resp = fmt.Sprintf("Set %s as default csv", splStr[1])
			}
			if strings.Contains(update.Message.Text, "/list_csv") {
				files, err := helpers.GetCSVFiles()
				if err != nil {
					fmt.Println(err)
				}
				resp = strings.Join(files, ", ")
			}
			if strings.Contains(update.Message.Text, "/add") {
				if defaultCsvName == "" {
					resp = "Set default csv filename first"
				} else {
					splStr := strings.Split(update.Message.Text, " ")
					err = helpers.InsertNewValueUnderHeader("internal/repository/"+defaultCsvName, splStr[1], splStr[2])
					if err != nil {
						fmt.Println(err)
					}
					resp = fmt.Sprintf("Added %s under %s", splStr[2], splStr[1])
				}
			}
			if strings.Contains(update.Message.Text, "/ya_file") {
				defaultExcelName := "tmp.xlsx"
				yacl.GetYDFileByPath(os.Getenv("YDFILE"), defaultExcelName)
				helpers.ConvertToCSV(defaultExcelName)
				resp = "File downloaded and converted to CSV"
			}
			if strings.Contains(update.Message.Text, "/ya_list") {
				list := yacl.GetYDList()
				resp = strings.Join(list, ", ")
			}
			if strings.Contains(update.Message.Text, "/ya_upload") {
				splStr := strings.Split(update.Message.Text, " ")
				if len(splStr) < 1 {
					resp = "Please specify filename"
				} else {
					err = yacl.SaveFileToYD(splStr[1])
					if err != nil {
						resp = "Error uploading file to Yandex Disk " + err.Error()
					}
					resp = "File uploaded to Yandex Disk"
				}
			}
			if strings.Contains(update.Message.Text, "/help") {
				resp = "/headers - get headers from default csv\n" +
					"/set_csv <filename> - set default csv\n" +
					"/list_csv - list all csv files\n" +
					"/add <header> <value> - add value under header\n" +
					"/ya_file - download file from Yandex Disk\n" +
					"/ya_list - list files from Yandex Disk\n" +
					"/ya_upload <filename> - upload file to Yandex Disk\n"
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
