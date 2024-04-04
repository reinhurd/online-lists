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

var lastChatID int64
var lastMessageID int

var defaultCsvName string

type TGBot struct {
	yacl *yandex.Client
	bot  *tgbotapi.BotAPI
	u    tgbotapi.UpdateConfig
}

func (t *TGBot) GetUpdatesChan() tgbotapi.UpdatesChannel {
	return t.bot.GetUpdatesChan(t.u)
}

func (t *TGBot) Send(chatID int64, messageID int, resp string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, resp)
	msg.ReplyToMessageID = messageID
	return t.bot.Send(msg)
}

func (t *TGBot) SendToLastChat(resp string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(lastChatID, resp)
	msg.ReplyToMessageID = lastMessageID
	return t.bot.Send(msg)
}

func (t *TGBot) HandleUpdate(updates tgbotapi.UpdatesChannel) error {
	var err error

	for update := range updates {
		var resp string
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			resp = "Hello, " + update.Message.From.UserName + "!" + " You said: " + update.Message.Text

			switch {
			case update.Message.Text == "/headers":
				//todo fix if defaultCsvName is empty
				res := helpers.GetCSVHeaders("internal/repository/" + defaultCsvName)
				resp = strings.Join(res, ", ")
			case strings.Contains(update.Message.Text, "/set_csv"):
				splStr := strings.Split(update.Message.Text, " ")
				defaultCsvName = splStr[1]
				resp = fmt.Sprintf("Set %s as default csv", splStr[1])
			case strings.Contains(update.Message.Text, "/list_csv"):
				files, err := helpers.GetCSVFiles()
				if err != nil {
					fmt.Println(err)
				}
				resp = strings.Join(files, ", ")
			case strings.Contains(update.Message.Text, "/add"):
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
			case strings.Contains(update.Message.Text, "/ya_file"):
				defaultExcelName := "tmp.xlsx"
				t.yacl.GetYDFileByPath(os.Getenv("YDFILE"), defaultExcelName)
				helpers.ConvertToCSV(defaultExcelName)
				resp = "File downloaded and converted to CSV"
			case strings.Contains(update.Message.Text, "/ya_list"):
				list := t.yacl.GetYDList()
				resp = strings.Join(list, ", ")
			case strings.Contains(update.Message.Text, "/ya_upload"):
				splStr := strings.Split(update.Message.Text, " ")
				if len(splStr) < 1 {
					resp = "Please specify filename"
				} else {
					err = t.yacl.SaveFileToYD(splStr[1])
					if err != nil {
						resp = "Error uploading file to Yandex Disk " + err.Error()
					}
					resp = "File uploaded to Yandex Disk"
				}
			case strings.Contains(update.Message.Text, "/help"):
				resp = "/headers - get headers from default csv\n" +
					"/set_csv <filename> - set default csv\n" +
					"/list_csv - list all csv files\n" +
					"/add <header> <value> - add value under header\n" +
					"/ya_file - download file from Yandex Disk\n" +
					"/ya_list - list files from Yandex Disk\n" +
					"/ya_upload <filename> - upload file to Yandex Disk\n"
			}

			lastChatID = update.Message.Chat.ID
			lastMessageID = update.Message.MessageID

			_, err = t.Send(update.Message.Chat.ID, update.Message.MessageID, resp)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func StartBot(tgToken string, yacl *yandex.Client, isDebug bool) (*TGBot, error) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = isDebug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	tgbot := &TGBot{
		yacl: yacl,
		bot:  bot,
		u:    u,
	}

	return tgbot, nil
}
