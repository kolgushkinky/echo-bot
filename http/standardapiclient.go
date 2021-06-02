package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type stdTelegramAPIClient struct {
	baseUrl string
}

func decodeResponseOnSuccess(body io.Reader, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}

func (c *stdTelegramAPIClient) Get(command string, v interface{}) error {
	url := c.baseUrl + command
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponseOnSuccess(resp.Body, v)
}

func (c *stdTelegramAPIClient) Post(command string, data url.Values, v interface{}) error {
	url := c.baseUrl + command
	resp, err := http.PostForm(url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponseOnSuccess(resp.Body, v)
}

func NewStdTelegramAPIClient(token string) (TelegramAPIClient, error) {
	if len(token) == 0 {
		return nil, errors.New("fatal: token cannot be empty")
	}

	return &stdTelegramAPIClient{baseUrl: "https://api.telegram.org/bot" + token + "/"}, nil
}
