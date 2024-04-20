package postgres

import (
	"time"

	"github.com/nawaltni/tracker/domain"
	"github.com/paulmach/orb"
	"gorm.io/datatypes"
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
		PhoneMeta: PhoneMeta{
			DeviceID:   in.PhoneMeta.DeviceID,
			Brand:      in.PhoneMeta.Brand,
			Model:      in.PhoneMeta.Model,
			OS:         in.PhoneMeta.OS,
			AppVersion: in.PhoneMeta.AppVersion,
			Carrier:    in.PhoneMeta.Carrier,
			Battery:    in.PhoneMeta.Battery,
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
		PhoneMeta: domain.PhoneMeta{
			DeviceID:   in.PhoneMeta.DeviceID,
			Brand:      in.PhoneMeta.Brand,
			Model:      in.PhoneMeta.Model,
			OS:         in.PhoneMeta.OS,
			AppVersion: in.PhoneMeta.AppVersion,
			Carrier:    in.PhoneMeta.Carrier,
			Battery:    in.PhoneMeta.Battery,
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

// ToModelTimeTrackingSession converts a domain TimeTrackingSession to a postgres TimeTrackingSession.
func ToModelTimeTrackingSession(in *domain.TimeTrackingSession) TimeTrackingSession {
	// Convert breaks
	breaks := make([]TimeTrackingBreak, len(in.Breaks))
	for i, b := range in.Breaks {
		breaks[i] = TimeTrackingBreak{
			BreakID:   b.BreakID,
			SessionID: in.SessionID,
			StartTime: b.StartTime,
			EndTime:   b.EndTime,
		}
	}
	return TimeTrackingSession{
		SessionID:      in.SessionID,
		UserID:         in.UserID,
		Status:         TimeTrackingStatus(in.Status),
		CheckedInTime:  in.CheckedInTime,
		CheckedOutTime: in.CheckedOutTime,
		TotalWorkTime:  datatypes.Time(in.TotalWorkTime),
		TotalBreakTime: datatypes.Time(in.TotalBreakTime),
		LastKnownLocation: GeoPoint{
			Point: orb.Point{float64(in.LastKnownLocation.Longitude), float64(in.LastKnownLocation.Latitude)},
		},
		Breaks: breaks,
	}
}

// ToDomainTimeTrackingSession converts a postgres TimeTrackingSession to a domain TimeTrackingSession.
func ToDomainTimeTrackingSession(in TimeTrackingSession) *domain.TimeTrackingSession {
	domainSession := &domain.TimeTrackingSession{
		SessionID:      in.SessionID,
		UserID:         in.UserID,
		Status:         domain.TimeTrackingStatus(in.Status),
		CheckedInTime:  in.CheckedInTime,
		CheckedOutTime: in.CheckedOutTime,
		TotalWorkTime:  time.Duration(in.TotalWorkTime),
		TotalBreakTime: time.Duration(in.TotalBreakTime),
		LastKnownLocation: domain.GeoPoint{
			Latitude:  float32(in.LastKnownLocation.Point.Y()),
			Longitude: float32(in.LastKnownLocation.Point.X()),
		},
	}
	// Convert breaks
	domainSession.Breaks = make([]domain.TimeTrackingBreak, len(in.Breaks))
	for i, b := range in.Breaks {
		domainSession.Breaks[i] = domain.TimeTrackingBreak{
			BreakID:   b.BreakID,
			StartTime: b.StartTime,
			EndTime:   b.EndTime,
		}
	}
	return domainSession
}

// // ToModelTimeTrackingBreak converts a domain TimeTrackingBreak to a postgres TimeTrackingBreak.
// func ToModelTimeTrackingBreak(in domain.TimeTrackingBreak) TimeTrackingBreak {
// 	return TimeTrackingBreak{
// 		BreakID:   in.BreakID,
// 		SessionID: in.SessionID,
// 		StartTime: in.StartTime,
// 		EndTime:   in.EndTime,
// 	}
// }

// // ToDomainTimeTrackingBreak converts a postgres TimeTrackingBreak to a domain TimeTrackingBreak.
// func ToDomainTimeTrackingBreak(in TimeTrackingBreak) domain.TimeTrackingBreak {
// 	return domain.TimeTrackingBreak{
// 		BreakID:   in.BreakID,
// 		SessionID: in.SessionID,
// 		StartTime: in.StartTime,
// 		EndTime:   in.EndTime,
// 	}
// }
