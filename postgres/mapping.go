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
		UpdatedAt:  in.UpdatedAt,
		PlaceID:    in.PlaceID,
		PlaceName:  in.PlaceName,
		CheckedIn:  in.CheckedIn,
		CheckedOut: in.CheckedOut,
		PhoneMeta: PhoneMeta{
			DeviceID:   in.PhoneMeta.DeviceID,
			Brand:      in.PhoneMeta.Brand,
			Model:      in.PhoneMeta.Model,
			OS:         in.PhoneMeta.OS,
			AppVersion: in.PhoneMeta.AppVersion,
			Carrier:    in.PhoneMeta.Carrier,
			Battery:    in.PhoneMeta.Battery,
		},
		Reference:     in.Reference,
		BackendUserID: in.BackendUserID,
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
		UpdatedAt: in.UpdatedAt,
		PlaceID:   in.PlaceID,
		PlaceName: in.PlaceName,
		PhoneMeta: domain.PhoneMeta{
			DeviceID:   in.PhoneMeta.DeviceID,
			Brand:      in.PhoneMeta.Brand,
			Model:      in.PhoneMeta.Model,
			OS:         in.PhoneMeta.OS,
			AppVersion: in.PhoneMeta.AppVersion,
			Carrier:    in.PhoneMeta.Carrier,
			Battery:    in.PhoneMeta.Battery,
		},
		Reference:     in.Reference,
		BackendUserID: in.BackendUserID,
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
