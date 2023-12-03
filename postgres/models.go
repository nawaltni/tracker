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
	UserID     string `gorm:"primary_key"`
	Reference  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Latitude   float32
	Longitude  float32
	PlaceID    *string
	PlaceName  *string
	CheckedIn  *time.Time
	CheckedOut *time.Time
	Location   GeoPoint
	PhoneMeta  PhoneMeta `gorm:"foreignKey:user_id"`
}

type PhoneMeta struct {
	UserID     string `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ID         string
	DeviceID   string
	Brand      string
	Model      string
	OS         string
	AppVersion string
	Carrier    string
	Battery    int
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

	var in []byte
	switch v := input.(type) {
	case []byte:
		in = v
	case string:
		in = []byte(v)
	default:
		return fmt.Errorf("invalid type for GeoPoint: %v", v)
	}
	var p orb.Point
	err := ewkb.Scanner(&p).Scan(in)
	if err != nil {
		return err
	}
	g.Point = p

	return nil
}
