package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

// получение сообщений (updates)
func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	//do request
	data, err := c.doRequest(getUpdatesMethod, q)

	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("can't unmarshal json: %w", err)
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
		return nil, fmt.Errorf("can't do request: %w", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("can't recieve response: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("can't read info from response: %w", err)
	}

	return body, err

}

func (c *Client) SendMessage(chatId int, text string) error {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, sendMessageMethod),
	}
	log.Println(u.String())

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)
	q.Add("parse_mode", "Markdown")
	//buf, err := json.Marshal(q)
	//if err != nil {
	//	return err
	//}

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	req.Header.Set("Content-Type", "application/json")

	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)

	if err != nil {
		return fmt.Errorf("can't recieve response: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	_, err = io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("can't read info from response: %w", err)
	}

	return err

}
