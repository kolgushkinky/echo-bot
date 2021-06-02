package bot

import "github.com/desotech-it/telegram-echo-bot/http"

type Bot struct {
	httpClient http.TelegramAPIClient
}

func NewBot(httpClient http.TelegramAPIClient) *Bot {
	return &Bot{httpClient: httpClient}
}
