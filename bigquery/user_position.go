package bigquery

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type UserPosition struct {
	UUID       string                 `json:"uuid"`
	CreatedAt  time.Time              `json:"created_at"`
	UserID     string                 `json:"user_id"`
	Location   string                 `json:"location"`
	PlaceID    bigquery.NullString    `json:"place_id"`
	PlaceName  bigquery.NullString    `json:"place_name"`
	CheckedIn  bigquery.NullTimestamp `json:"checked_in"`
	CheckecOut bigquery.NullTimestamp `json:"checked_out"`
	PhoneMeta  *PhoneMeta             `json:"phone_meta"`
}

type PhoneMeta struct {
	ID         string `json:"id"`
	DeviceID   string `json:"device_id"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	OS         string `json:"os"`
	AppVersion string `json:"app_version"`
	Carrier    string `json:"carrier"`
	Battery    int    `json:"battery"`
}
