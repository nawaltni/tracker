package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserPosition struct {
	ID            string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID        string     `gorm:"type:uuid;not null"`
	Latitude      float64    `gorm:"type:double precision;not null"`
	Longitude     float64    `gorm:"type:double precision;not null"`
	CreatedAt     time.Time  `gorm:"type:timestamp with time zone;not null"`
	PlaceID       *string    `gorm:"type:uuid"`
	CheckedIn     *time.Time `gorm:"type:timestamp with time zone"`
	CheckedOut    *time.Time `gorm:"type:timestamp with time zone"`
	Location      GeoPoint
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

type GeoPoint struct {
	Point orb.Point
}

func (g GeoPoint) GormDataType() string {
	return "geometry(Point, 4326)"
}

func (g GeoPoint) GormDBDataType() string {
	return "geometry(Point, 4326)"
}
func (g GeoPoint) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {

	return clause.Expr{
		SQL:  "ST_GeomFromEWKB(?)",
		Vars: []interface{}{ewkb.Value(g.Point, 4326)},
	}
}

func (g *GeoPoint) Scan(input interface{}) error {
	fmt.Printf("Scanning GeoPoint: %v\n", input)
	var p orb.Point
	err := ewkb.Scanner(&input).Scan(&p)
	if err != nil {
		return err
	}
	g.Point = p

	return nil
}
