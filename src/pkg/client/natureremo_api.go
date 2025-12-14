package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kenta-ja8/home-k8s-app/pkg/entity"
	"github.com/kenta-ja8/home-k8s-app/pkg/logger"
	"github.com/pkg/errors"
)

type NatureremoClient struct {
	cfg *entity.Config
}

func NewNatureremoClient(cfg *entity.Config) *NatureremoClient {
	return &NatureremoClient{
		cfg: cfg,
	}
}

const devicesEndpoint = "https://api.nature.global/1/devices"

func (c *NatureremoClient) GetDevices(
	ctx context.Context,
) ([]entity.Device, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, devicesEndpoint, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.NATUREREMO_TOKEN)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.New(
			fmt.Sprintf("Nature Remo API request failed with status: %s, body: %s", resp.Status, string(b)),
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	logger.Info("Nature Remo API response body: %s", string(body))

	var devices []entity.Device
	if err := json.Unmarshal(body, &devices); err != nil {
		return nil, errors.WithStack(err)
	}

	return devices, nil
}
