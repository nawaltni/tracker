package bigquery

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type UserPosition struct {
	UUID       string                 `bigquery:"uuid"`
	CreatedAt  time.Time              `bigquery:"created_at"`
	UserID     string                 `bigquery:"user_id"`
	Name       string                 `bigquery:"name"`
	Location   string                 `bigquery:"location"`
	PlaceID    bigquery.NullString    `bigquery:"place_id"`
	PlaceName  bigquery.NullString    `bigquery:"place_name"`
	CheckedIn  bigquery.NullTimestamp `bigquery:"checked_in"`
	CheckecOut bigquery.NullTimestamp `bigquery:"checked_out"`
	PhoneMeta  *PhoneMeta             `bigquery:"phone_meta"`
}

type PhoneMeta struct {
	DeviceID   string `bigquery:"device_id"`
	Brand      string `bigquery:"brand"`
	Model      string `bigquery:"model"`
	OS         string `bigquery:"os"`
	AppVersion string `bigquery:"app_version"`
	Carrier    string `bigquery:"carrier"`
	Battery    int    `bigquery:"battery"`
}
