package postgres

import (
	"github.com/nawaltni/tracker/domain"

	"github.com/paulmach/orb"
)

// ToModelUserPosition converts a domain UserPosition to a postgres UserPosition.
func ToModelUserPosition(in *domain.UserPosition) UserPosition {
	up := UserPosition{
		UserID:       in.UserID,
		Latitude:     in.Location.Latitude,
		Longitude:    in.Location.Longitude,
		Timestamp:    in.CreatedAt,
		PlaceID:      in.PlaceID,
		CheckedInAt:  in.CheckedInAt,
		CheckedOutAt: in.CheckedOutAt,
		Location:     orb.Point{in.Location.Longitude, in.Location.Latitude},
		PhoneMetadata: PhoneMetadata{
			DeviceID:    in.Metadata.DeviceID,
			Model:       in.Metadata.Model,
			OSVersion:   in.Metadata.OSVersion,
			Carrier:     in.Metadata.Carrier,
			CorporateID: in.Metadata.CorporateID,
		},
	}
	return up
}

// ToDomainUserPosition converts a postgres UserPosition to a domain UserPosition.
func ToDomainUserPosition(in *UserPosition) *domain.UserPosition {
	domainUP := &domain.UserPosition{
		UserID:    in.UserID,
		Location:  domain.GeoPoint{Latitude: in.Latitude, Longitude: in.Longitude},
		CreatedAt: in.Timestamp,
		PlaceID:   in.PlaceID,
		Metadata: domain.PhoneMetadata{
			DeviceID:    in.PhoneMetadata.DeviceID,
			Model:       in.PhoneMetadata.Model,
			OSVersion:   in.PhoneMetadata.OSVersion,
			Carrier:     in.PhoneMetadata.Carrier,
			CorporateID: in.PhoneMetadata.CorporateID,
		},
	}
	// Conditional assignment for optional fields
	if in.CheckedInAt != nil {
		domainUP.CheckedInAt = in.CheckedInAt
	}
	if in.CheckedOutAt != nil {
		domainUP.CheckedOutAt = in.CheckedOutAt
	}
	// If you have a PlaceName field in the model, you would convert it here as well.
	return domainUP
}
