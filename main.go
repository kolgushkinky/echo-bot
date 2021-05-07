package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nervigalilei/echo-bot/bot"
	"github.com/nervigalilei/echo-bot/http"
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
			response := fmt.Sprintf("Ciao %s! Hai scritto: %s", userFirstName, text)
			bot.SendMessage(userId, response)
			offset = update.ID + 1
		}
	}
}
