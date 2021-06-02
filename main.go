package main

import (
	"log"
	"os"

	"github.com/desotech-it/telegram-echo-bot/bot"
	"github.com/desotech-it/telegram-echo-bot/http"
)

func abortOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	client, err := http.NewStdTelegramAPIClient(os.Getenv("TG_TOKEN"))
	abortOnError(err)
	bot := bot.NewBot(client)
	offset := int64(0)
	for {
		updates, err := bot.GetUpdates(offset)
		abortOnError(err)
		for _, update := range updates {
			text := update.Message.Text
			userFirstName := update.Message.From.FirstName
			userId := update.Message.Chat.ID
			response := "Hi " + userFirstName + "! You wrote: " + text
			bot.SendMessage(userId, response)
			offset = update.ID + 1
		}
	}
}
