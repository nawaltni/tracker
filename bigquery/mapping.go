package bigquery

import (
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
