package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/lysenkopavlo/telegram-bot/internal/helpers/e"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string //
	client   http.Client
}

func (c *Client) SendMessage(chatId int, text string) error {
	params := url.Values{}
	params.Add("chat_id", strconv.Itoa(chatId))
	params.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, params)
	if err != nil {
		return e.WrapError("Error while sending message: %w", err)
	}

	return nil
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	params := url.Values{}
	params.Add("offset", strconv.Itoa(offset))
	params.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, params)

	if err != nil {
		return nil, e.WrapError("Error while NewRequest: %w", err)

	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.WrapError("Error while NewRequest: %w", err)

	}
	return res.Result, nil

}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, e.WrapError("Error while NewRequest: %w", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.WrapError("Error while NewRequest: %w", err)

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, e.WrapError("Error while NewRequest: %w", err)
	}
	return body, nil
}

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}
