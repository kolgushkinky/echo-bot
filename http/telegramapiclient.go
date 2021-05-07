package http

import "net/url"

type TelegramAPIClient interface {
	Get(command string, v interface{}) error
	Post(command string, data url.Values, v interface{}) error
}
