package grpc

import (
	"context"

	pb "github.com/nawaltni/api/gen/go/nawalt/tracker/v1"
	"github.com/nawaltni/tracker/domain"
	"github.com/nawaltni/tracker/services"
	"github.com/pkg/errors"
)

// NewServer creates a new TrackingService server.
func NewServer(srvs *services.Services) *Server {
	return &Server{
		services: srvs,
	}
}

// RecordPosition records a user's location update.
func (s *Server) RecordPosition(ctx context.Context, req *pb.RecordPositionRequest) (*pb.RecordPositionResponse, error) {
	err := s.services.UserPositionService.RecordPosition(
		ctx,
		req.UserId,
		domain.GeoPoint{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		req.Timestamp.AsTime(),
		req.ClientId,
		domain.PhoneMetadata{
			DeviceID:    req.Metadata.DeviceId,
			Model:       req.Metadata.Model,
			OSVersion:   req.Metadata.OsVersion,
			Carrier:     req.Metadata.Carrier,
			CorporateID: req.Metadata.CorporateId,
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to record position")
	}

	return &pb.RecordPositionResponse{
		Success: true,
		Message: "Position recorded successfully",
	}, nil
}

// GetUserPosition retrieves a user's current position.
func (s *Server) GetUserPosition(ctx context.Context, req *pb.GetUserPositionRequest) (*pb.GetUserPositionResponse, error) {
	userPosition, err := s.services.UserPositionService.GetUserPosition(ctx, req.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user position")
	}

	return &pb.GetUserPositionResponse{
		UserPosition: convertToPbUserPosition(userPosition),
	}, nil
}

// GetUsersPositionByCoordinates retrieves a user's position by coordinates.
func (s *Server) GetUsersPositionByCoordinates(ctx context.Context, req *pb.GetUsersPositionByCoordinatesRequest) (*pb.GetUsersPositionByCoordinatesResponse, error) {
	userPositions, err := s.services.UserPositionService.GetUsersPositionByCoordinates(ctx, req.Latitude, req.Longitude, int(req.Distance))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users' positions by coordinates")
	}

	// Convert domain.UserPosition slice to pb.UserPosition slice
	var pbUserPositions []*pb.UserPosition
	for _, up := range userPositions {
		pbUserPositions = append(pbUserPositions, convertToPbUserPosition(&up))
	}

	return &pb.GetUsersPositionByCoordinatesResponse{
		UserPositions: pbUserPositions,
	}, nil
}
