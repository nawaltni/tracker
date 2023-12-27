package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/nawaltni/tracker/domain"
)

// Client is a bigquery client wrapper we use to communicate with bigquery
type Client struct {
	client *bigquery.Client
}

// NewClient returns a new bigquery client
func NewClient(project string) (*Client, error) {
	c, err := bigquery.NewClient(context.Background(), project)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: c,
	}, nil
}

func (c *Client) RecordUserPosition(ctx context.Context, userPosition domain.UserPosition) error {

	data := []UserPosition{
		mapToUserPosition(userPosition),
	}

	inserter := c.client.Dataset("development").Table("users-position").Inserter()

	if err := inserter.Put(ctx, data); err != nil {
		return fmt.Errorf("error inserting rows: %w", err)
	}
	return nil
}
