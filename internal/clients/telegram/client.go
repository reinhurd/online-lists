package telegram

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"online-lists/internal/service"
)

var lastChatID int64

type TGBot struct {
	olSvc *service.Service
	bot   *tgbotapi.BotAPI
	u     tgbotapi.UpdateConfig
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
	msg.ReplyToMessageID = 0
	return t.bot.Send(msg)
}

func (t *TGBot) HandleUpdate(updates tgbotapi.UpdatesChannel) error {
	var err error

	for update := range updates {
		var resp string
		if update.Message != nil { // If we got a message
			log.Info().Msgf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			resp = "Hello, " + update.Message.From.UserName + "!" + " You said: " + update.Message.Text

			switch {
			case update.Message.Text == "/headers":
				resp = t.olSvc.GetHeaders()
			case strings.Contains(update.Message.Text, "/set_csv"):
				splStr := strings.Split(update.Message.Text, " ")
				resp = t.olSvc.SetDefaultCsv(splStr[1])
			case strings.Contains(update.Message.Text, "/list_csv"):
				resp = t.olSvc.ListCsv()
			case strings.Contains(update.Message.Text, "/add"):
				splStr := strings.Split(update.Message.Text, " ")
				resp = t.olSvc.Add(splStr[1], splStr[2])
			case strings.Contains(update.Message.Text, "/ya_file"):
				var filename string
				splStr := strings.Split(update.Message.Text, " ")
				if len(splStr) > 1 {
					filename = splStr[1]
				}
				resp = t.olSvc.YAFile(filename)
			case strings.Contains(update.Message.Text, "/ya_list"):
				resp = strings.Join(t.olSvc.GetYaList(), ", ")
			case strings.Contains(update.Message.Text, "/ya_upload"):
				splStr := strings.Split(update.Message.Text, " ")
				if len(splStr) < 2 {
					resp = "Please specify filename"
				} else {
					resp = t.olSvc.YAUpload(splStr[1])
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

			_, err = t.Send(update.Message.Chat.ID, update.Message.MessageID, resp)
			if err != nil {
				log.Err(err).Msg("send error")
			}
		}
	}
	return nil
}

func StartBot(tgToken string, olSvc *service.Service, isDebug bool) (*TGBot, error) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Fatal().Err(err).Msg("tgbotapi.NewBotAPI doesn't start")
	}

	bot.Debug = isDebug

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	tgbot := &TGBot{
		olSvc: olSvc,
		bot:   bot,
		u:     u,
	}

	return tgbot, nil
}
