package bigquery

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/nawaltni/tracker/domain"
)

func mapToUserPosition(in domain.UserPosition) UserPosition {
	out := UserPosition{
		UUID:      in.Reference,
		CreatedAt: in.CreatedAt,
		UserID:    in.UserID,
		Location:  in.Location.String(),
	}

	if in.PlaceID != nil {
		out.PlaceID = bigquery.NullString{StringVal: *in.PlaceID, Valid: true}
	}

	if in.PlaceName != nil {
		out.PlaceName = bigquery.NullString{StringVal: *in.PlaceName, Valid: true}
	}

	if in.CheckedIn != nil {
		out.CheckedIn = bigquery.NullTimestamp{Timestamp: *in.CheckedIn, Valid: true}
	}

	if in.CheckedOut != nil {
		out.CheckecOut = bigquery.NullTimestamp{Timestamp: *in.CheckedOut, Valid: true}
	}

	out.PhoneMeta = mapToPhoneMeta(in.PhoneMeta)

	return out
}

func mapToPhoneMeta(in domain.PhoneMeta) *PhoneMeta {
	out := &PhoneMeta{
		DeviceID:   in.DeviceID,
		Brand:      in.Brand,
		Model:      in.Model,
		OS:         in.OS,
		AppVersion: in.AppVersion,
		Battery:    in.Battery,
	}

	return out
}

func mapToDomainPhoneMeta(in PhoneMeta) domain.PhoneMeta {
	out := domain.PhoneMeta{
		DeviceID:   in.DeviceID,
		Brand:      in.Brand,
		Model:      in.Model,
		OS:         in.OS,
		AppVersion: in.AppVersion,
		Battery:    in.Battery,
	}

	return out
}

func mapDomainLocationFromString(in string) domain.GeoPoint {
	// in = "POINT(0.000000 0.000000)"
	var lat, long float32
	fmt.Sscanf(in, "POINT(%f %f)", &long, &lat)
	return domain.GeoPoint{
		Latitude:  lat,
		Longitude: long,
	}
}

func mapToDomainUserPosition(in UserPosition) domain.UserPosition {
	out := domain.UserPosition{
		Reference: in.UUID,
		CreatedAt: in.CreatedAt,
		UserID:    in.UserID,
		Location:  mapDomainLocationFromString(in.Location),
	}

	if in.PlaceID.Valid {
		out.PlaceID = &in.PlaceID.StringVal
	}

	if in.PlaceName.Valid {
		out.PlaceName = &in.PlaceName.StringVal
	}

	if in.CheckedIn.Valid {
		out.CheckedIn = &in.CheckedIn.Timestamp
	}

	if in.CheckecOut.Valid {
		out.CheckedOut = &in.CheckecOut.Timestamp
	}

	if in.PhoneMeta != nil {
		out.PhoneMeta = mapToDomainPhoneMeta(*in.PhoneMeta)
	}

	return out
}
