package main

import (
	"log"
	"time"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	greet := config.GreetingMessage

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	grouped := false
	media := []interface{}{}
	
	inactivityTimer := time.NewTimer(time.Duration(u.Timeout/4) * time.Second)

	for {
		select {
		case update := <-updates:
			privateChatID := update.Message.Chat.ID
			text := update.Message.Text
			if update.Message != nil && update.Message.ReplyToMessage == nil && update.Message.Chat.Type == "private" {
				if update.Message.MediaGroupID != "" {
					grouped = true
					
					switch {
					case update.Message.Photo != nil :
						photo := update.Message.Photo
						fileID := photo[len(photo)-1].FileID

						mediaItem := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(fileID))
						if update.Message.Caption != "" {
							mediaItem.Caption = update.Message.Caption
						}
						media = append(media, mediaItem)
					
					case update.Message.Video != nil :
						videoItem := *update.Message.Video
						fileID := videoItem.FileID
						
						mediaItem := tgbotapi.NewInputMediaVideo(tgbotapi.FileID(fileID))
						if update.Message.Caption != "" {
							mediaItem.Caption = update.Message.Caption
						}
						media = append(media, mediaItem)
					default:
						copyMessageConfig := tgbotapi.NewCopyMessage(chatID, update.Message.Chat.ID, update.Message.MessageID)
						_, err := bot.CopyMessage(copyMessageConfig)
						if err != nil {
							log.Println(err)
						}
					}
					
				} else if text == "/start"{
					greetingMsg := tgbotapi.NewMessage(privateChatID, greet)
					_, err := bot.Send(greetingMsg)
					if err != nil {
						log.Println(err)
					}
				} else {
					copyMessageConfig := tgbotapi.NewCopyMessage(chatID, update.Message.Chat.ID, update.Message.MessageID)
					_, err := bot.CopyMessage(copyMessageConfig)
					if err != nil {
						log.Println(err)
					}
				}
			}

		case <-inactivityTimer.C:
			if grouped {
				grouped = false
				config := tgbotapi.NewMediaGroup(chatID, media)
				_, err := bot.SendMediaGroup(config)
				if err != nil {
					log.Println(err)
				}
				media = []interface{}{}
			}
			inactivityTimer.Reset(time.Duration(u.Timeout/4) * time.Second)
		}
	}
}