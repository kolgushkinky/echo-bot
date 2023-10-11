package main

import (
	"log"
	"time"
	//"net/http"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"image"
	//"image/jpeg"
	//"bytes"
	"os"
    "encoding/json"
)

type Config struct {
    BotToken       string `json:"bot_token"`
    ChatID         int64  `json:"chat_id"`
    GreetingMessage string `json:"greeting_message"`
}

func main() {

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var config Config

	if err := decoder.Decode(&config); err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	// where to resend messages
	var chatID int64 = config.ChatID

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//currentMediaGroupID := ""
	grouped := false
	media := []interface{}{}
	// Initialize a timer for checking inactivity
	inactivityTimer := time.NewTimer(time.Duration(u.Timeout/4) * time.Second)

	for {
		select {
		case update := <-updates:
			if update.Message != nil && update.Message.ReplyToMessage == nil {
				if update.Message.MediaGroupID != "" {
					grouped = true
					log.Println("MediaGroupID not empty ", update.Message.MediaGroupID)

					// Construct media group from received media files
					photo := update.Message.Photo
					fileID := photo[len(photo)-1].FileID
					fileConfig := tgbotapi.FileConfig{FileID: fileID}
					file, err := bot.GetFile(fileConfig)
					if err != nil {
						log.Println(err)
						continue
					}

					fileURL := file.Link(bot.Token)
					log.Println(fileURL, "is file url")

					mediaItem := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(fileID))
					if update.Message.Caption != "" {
						mediaItem.Caption = update.Message.Caption
					}
					media = append(media, mediaItem)
				} else {
					copyMessageConfig := tgbotapi.NewCopyMessage(chatID, update.Message.Chat.ID, update.Message.MessageID)
					_, err := bot.CopyMessage(copyMessageConfig)
					if err != nil {
						log.Println(err)
					}
				}
			}

		case <-inactivityTimer.C:
			// Timer expired, send pending media group (if any)
			if grouped {
				grouped = false
				config := tgbotapi.NewMediaGroup(chatID, media)
				log.Println("sending media", media)
				// Send the media group
				_, err := bot.SendMediaGroup(config)
				if err != nil {
					log.Println(err)
				}
				media = []interface{}{}
			}
			// Reset the inactivity timer
			inactivityTimer.Reset(time.Duration(u.Timeout/4) * time.Second)
		}
	}
}