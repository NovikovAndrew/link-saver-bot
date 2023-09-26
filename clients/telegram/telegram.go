package telegram

import (
	e "bot-saver/package/error"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	telegramGetUpdateMethod   = "getUpdates"
	telegramSendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func (c *Client) Updates(offSet, limit int) ([]Update, error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offSet))
	query.Add("limit", strconv.Itoa(limit))

	data, err := c.request(telegramGetUpdateMethod, query)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Response, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatID))
	query.Add("text", text)

	_, err := c.request(telegramSendMessageMethod, query)

	if err != nil {
		return e.Wrap(fmt.Sprintf("can't send message, message: %s", telegramSendMessageMethod), err)
	}

	return nil
}

func (c *Client) request(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(query.Encode(), method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, e.Wrap("can't to create request", err)
	}

	req.URL.RawQuery = query.Encode()

	response, err := c.client.Do(req)

	if err != nil {
		return nil, e.Wrap("can't to execute request", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, e.Wrap("can't to read body from response", err)
	}

	return body, nil
}

func newBasePath(token string) string {
	return "bot" + token
}
