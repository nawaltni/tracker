package bigquery

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/nawaltni/tracker/domain"
	"google.golang.org/api/iterator"
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

// RecordUserPosition records a user's position in bigquery
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

// GetUserPositionsSince retrieves a user's positions from bigquery since a given time
func (c *Client) GetUserPositionsSince(ctx context.Context, userID string, t time.Time) ([]domain.UserPosition, error) {
	q := c.client.Query(fmt.Sprintf("SELECT * FROM `development.users-position` WHERE user_id = '%s' AND created_at >= '%s'", userID, t.Format("2006-01-02 15:04:05")))
	it, err := q.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user position: %w", err)
	}

	var userPositions []domain.UserPosition
	for {
		var up UserPosition
		err := it.Next(&up)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error retrieving user position: %w", err)
		}
		userPositions = append(userPositions, mapToDomainUserPosition(up))
	}

	return userPositions, nil
}

//
