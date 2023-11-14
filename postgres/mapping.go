package postgres

import (
	"github.com/nawaltni/tracker/domain"
	"github.com/paulmach/orb"
)

// ToModelUserPosition converts a domain UserPosition to a postgres UserPosition.
func ToModelUserPosition(in *domain.UserPosition) UserPosition {
	up := UserPosition{
		UserID:     in.UserID,
		Latitude:   in.Location.Latitude,
		Longitude:  in.Location.Longitude,
		CreatedAt:  in.CreatedAt,
		PlaceID:    in.PlaceID,
		PlaceName:  in.PlaceName,
		CheckedIn:  in.CheckedIn,
		CheckedOut: in.CheckedOut,
		PhoneMetadata: PhoneMetadata{
			DeviceID:    in.Metadata.DeviceID,
			Model:       in.Metadata.Model,
			OSVersion:   in.Metadata.OSVersion,
			Carrier:     in.Metadata.Carrier,
			CorporateID: in.Metadata.CorporateID,
		},
	}

	point := orb.Point{float64(in.Location.Longitude), float64(in.Location.Latitude)}

	// point := geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{in.Location.Longitude, in.Location.Latitude}).SetSRID(4326)
	up.Location = GeoPoint{Point: point}
	return up
}

// ToDomainUserPosition converts a postgres UserPosition to a domain UserPosition.
func ToDomainUserPosition(in UserPosition) *domain.UserPosition {
	domainUP := &domain.UserPosition{
		UserID:    in.UserID,
		Location:  domain.GeoPoint{Latitude: in.Latitude, Longitude: in.Longitude},
		CreatedAt: in.CreatedAt,
		PlaceID:   in.PlaceID,
		PlaceName: in.PlaceName,
		Metadata: domain.PhoneMetadata{
			DeviceID:    in.PhoneMetadata.DeviceID,
			Model:       in.PhoneMetadata.Model,
			OSVersion:   in.PhoneMetadata.OSVersion,
			Carrier:     in.PhoneMetadata.Carrier,
			CorporateID: in.PhoneMetadata.CorporateID,
		},
	}
	// Conditional assignment for optional fields
	if in.CheckedIn != nil {
		domainUP.CheckedIn = in.CheckedIn
	}
	if in.CheckedOut != nil {
		domainUP.CheckedOut = in.CheckedOut
	}
	// If you have a PlaceName field in the model, you would convert it here as well.
	return domainUP
}
