package pilot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultBaseURL = "https://tools.nexomechanics.com/api/pilot"

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

type Option func(*Client)

func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.Timeout = d }
}

func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

type SendResult struct {
	Remaining int
}

type PilotError struct {
	Message    string
	StatusCode int
}

func (e *PilotError) Error() string {
	return fmt.Sprintf("pilot: %s (status %d)", e.Message, e.StatusCode)
}

func (c *Client) Send(destination, message string) (*SendResult, error) {
	body, _ := json.Marshal(map[string]string{"destination": destination, "message": message})

	req, err := http.NewRequest("POST", c.baseURL+"/v1/forward", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]any
	json.NewDecoder(res.Body).Decode(&data)

	if res.StatusCode >= 400 {
		msg, _ := data["error"].(string)
		if msg == "" {
			msg = "request failed"
		}
		return nil, &PilotError{Message: msg, StatusCode: res.StatusCode}
	}

	remaining := 0
	if v, ok := data["remaining"].(float64); ok {
		remaining = int(v)
	}
	return &SendResult{Remaining: remaining}, nil
}
