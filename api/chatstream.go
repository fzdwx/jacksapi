package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type ChatStream struct {
	temperature     float64
	presencePenalty float64
	message         []ChatMessage
	c               *Client
}

func (c *ChatStream) Temperature(temperature float64) *ChatStream {
	c.temperature = temperature
	return c
}

func (c *ChatStream) PresencePenalty(presencePenalty float64) *ChatStream {
	c.presencePenalty = presencePenalty
	return c
}

func (c *ChatStream) Message(message []ChatMessage) *ChatStream {
	c.message = message
	return c
}

func (c *ChatStream) Model(model string) *ChatStream {
	c.c.Model = model
	return c
}

func (c *ChatStream) Do() (*http.Response, error) {
	client := http.DefaultClient
	body, err := c.buildBody()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, ApiBaseUrl+"/api/chat-stream", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Code", c.c.accessCode)
	req.Header.Set("Origin", ApiBaseUrl)
	req.Header.Set("Referer", ApiBaseUrl)
	req.Header.Set("Path", "v1/chat/completions")
	return client.Do(req)
}

func (c *ChatStream) DoWithCallback(callback Callback) {
	resp, err := c.Do()
	callback(resp, err)
}

func (c *ChatStream) buildBody() (io.Reader, error) {
	body := Body{
		Messages:        c.message,
		Stream:          true,
		Model:           c.c.Model,
		Temperature:     c.temperature,
		PresencePenalty: c.presencePenalty,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}
