package main

import (
	"log"
	"net/http"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"image"
	"image/jpeg"
	"bytes"
)

func main() {

	bot, err := tgbotapi.NewBotAPI("6153779033:AAH4uTyl22ftJ9ORIr880_CiZutQW6QtY3k")
	if err != nil {
		log.Fatal(err)
	}

	//where to resend messages
	var chatID int64 = -1878673641
	var channelUsername string = "@myBotTestGroupj"

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			text := update.Message.Text

			if update.Message.Photo != nil {
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
							ChannelUsername: channelUsername,
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
				_, err = bot.Send(tgbotapi.NewMessageToChannel(channelUsername, text))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
