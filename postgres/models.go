package postgres

import (
	"time"

	"github.com/paulmach/orb"
)

type UserPosition struct {
	ID            string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID        string     `gorm:"type:uuid;not null"`
	Latitude      float64    `gorm:"type:double precision;not null"`
	Longitude     float64    `gorm:"type:double precision;not null"`
	Timestamp     time.Time  `gorm:"type:timestamp with time zone;not null"`
	PlaceID       *string    `gorm:"type:uuid"`
	CheckedInAt   *time.Time `gorm:"type:timestamp with time zone"`
	CheckedOutAt  *time.Time `gorm:"type:timestamp with time zone"`
	Location      orb.Point  `gorm:"type:geometry(Point,4326)"`
	PhoneMetadata PhoneMetadata
}

type PhoneMetadata struct {
	UserPositionID string `gorm:"type:uuid;primary_key"`
	DeviceID       string `gorm:"type:varchar(255);not null"`
	Model          string `gorm:"type:varchar(255);not null"`
	OSVersion      string `gorm:"type:varchar(255);not null"`
	Carrier        string `gorm:"type:varchar(255);not null"`
	CorporateID    string `gorm:"type:varchar(255);not null"`
}
