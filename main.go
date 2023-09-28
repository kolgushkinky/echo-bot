package main

import (
	"log"
	"net/http"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"image"
	"image/jpeg"
	"bytes"
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

	//where to resend messages
	var chatID int64 = config.ChatID
	//greeted := make(map[int64]bool)
	greet := config.GreetingMessage

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			privateChatID := update.Message.Chat.ID
			text := update.Message.Text

			if update.Message.Sticker != nil {
				sticker := *update.Message.Sticker
				fileID := sticker.FileID
				
				msg := tgbotapi.NewStickerShare(chatID, fileID)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			} else if update.Message.Photo != nil {
				photo := *update.Message.Photo
				fileID := photo[len(photo)-1].FileID
				fileConfig := tgbotapi.FileConfig{FileID: fileID}
				file, err := bot.GetFile(fileConfig)
				if err != nil {
					log.Println(err)
					continue
				}

				fileURL := file.Link(bot.Token)
				response, err := http.Get(fileURL)
				if err != nil {
					log.Println(err)
					continue
				}
				defer response.Body.Close()

				img, _, err := image.Decode(response.Body)
				if err != nil {
					log.Println(err)
					continue
				}

				var buf bytes.Buffer
				err = jpeg.Encode(&buf, img, nil)
				if err != nil {
				log.Println(err)
				continue
				}

				msg := tgbotapi.PhotoConfig{
					BaseFile: tgbotapi.BaseFile{
						BaseChat:    tgbotapi.BaseChat{
							ChatID: chatID,
						},
						File: tgbotapi.FileReader{
							Name:   "image.jpg",
							Reader: &buf,
							Size:   int64(buf.Len()),
						},
					},
				}

				if update.Message.Caption != "" {
					comment := update.Message.Caption
					msg.Caption = comment		
				} 

				_, err = bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
				
			} else {
				if text != "/start" {
					_, err = bot.Send(tgbotapi.NewMessage(chatID, text))
					if err != nil {
						log.Println(err)
					}
				} else {
					greetingMsg := tgbotapi.NewMessage(privateChatID, greet)

					_, err := bot.Send(greetingMsg)
					if err != nil {
						log.Println(err)
					}
				}	
			}
		}
	}
}
