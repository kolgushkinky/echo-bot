package bot

import (
	"net/url"
	"strconv"
)

type getUpdatesResponse struct {
	Ok     bool     `json:"bool"`
	Result []Update `json:"result"`
}

func (b *Bot) GetUpdates(offset int64) ([]Update, error) {
	data := url.Values{}
	data.Add("offset", strconv.FormatInt(offset, 10))
	// TODO: Handle all kinds of updates
	data.Add("allowed_updates", "message")
	var apiResponse getUpdatesResponse
	err := b.httpClient.Post("getUpdates", data, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse.Result, nil
}
