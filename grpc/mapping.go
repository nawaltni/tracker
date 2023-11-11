package grpc

import (
	"github.com/nawaltni/tracker/domain"
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
