package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/kenta-ja8/home-k8s-app/pkg/entity"
)

type LineMessageClient struct {
	cfg *entity.Config
}

func NewLineMessageClient(cfg *entity.Config) *LineMessageClient {
	return &LineMessageClient{
		cfg: cfg,
	}
}

type PushPayload struct {
	To       string              `json:"to"`
	Messages []map[string]string `json:"messages"`
}

func (c *LineMessageClient) SendMessage(message string) error {

	payload := PushPayload{
		To: c.cfg.LINE_GROUP_ID,
		Messages: []map[string]string{
			{
				"type": "text",
				"text": message,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}

	req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/push", bytes.NewBuffer(body))
	if err != nil {
		return errors.WithStack(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.LINE_CHANNEL_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to send message to LINE. status: " + resp.Status)
	}

	return nil
}
