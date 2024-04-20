package grpc

import (
	"github.com/nawaltni/tracker/domain"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	common "github.com/nawaltni/api/gen/go/nawalt/common/v1"
	pb "github.com/nawaltni/api/gen/go/nawalt/tracker/v1"
)

// convertToPbUserPosition converts domain.UserPosition to pb.UserPosition
func convertToPbUserPosition(up *domain.UserPosition) *pb.UserPosition {
	// Implement conversion logic here
	return &pb.UserPosition{
		UserId:    up.UserID,
		Location:  &common.GeoPoint{Latitude: up.Location.Latitude, Longitude: up.Location.Longitude},
		Timestamp: timestamppb.New(up.CreatedAt),
		// Fill in other fields as necessary
	}
}

func convertSessionToProto(session domain.TimeTrackingSession) *pb.TimeTrackingInfo {
	info := &pb.TimeTrackingInfo{
		SessionId:         session.SessionID,
		UserId:            session.UserID,
		LastKnownLocation: &common.GeoPoint{Latitude: session.LastKnownLocation.Latitude, Longitude: session.LastKnownLocation.Longitude},
		Status:            pb.TimeTrackingStatus(session.Status),
		CheckinTime:       timestamppb.New(session.CheckedInTime),
		TotalWorkTime:     durationpb.New(session.TotalWorkTime),
		TotalBreakTime:    durationpb.New(session.TotalBreakTime),
	}

	if session.CheckedOutTime != nil {
		info.CheckoutTime = timestamppb.New(*session.CheckedOutTime)
	}

	info.LastKnownLocation = &common.GeoPoint{
		Latitude:  session.LastKnownLocation.Latitude,
		Longitude: session.LastKnownLocation.Longitude,
	}

	if len(session.Breaks) > 0 {
		info.Breaks = make([]*pb.TimeTrackingBreak, len(session.Breaks))
		for i, b := range session.Breaks {
			info.Breaks[i] = &pb.TimeTrackingBreak{
				StartTime: timestamppb.New(b.StartTime),
			}
			if b.EndTime != nil {
				info.Breaks[i].EndTime = timestamppb.New(*b.EndTime)
			}
		}
	}

	return info
}
