package bigquery

import (
	"context"
	"log/slog"

	"github.com/nawaltni/tracker/domain"
)

type Client struct{}

func NewClient() (*Client, error) {
	return &Client{}, nil
}

func (c *Client) RecordUserPosition(ctx context.Context, userPosition domain.UserPosition) error {
	slog.Info("Bigquery would record the user position", "info", userPosition)
	return nil
}
