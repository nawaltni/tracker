package grpc

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	pb "github.com/nawaltni/api/gen/go/nawalt/tracker/v1"
	"github.com/nawaltni/tracker/domain"
	"github.com/nawaltni/tracker/services"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewServer creates a new TrackingService server.
func NewServer(srvs *services.Services) *Server {
	return &Server{
		services: srvs,
	}
}

// RecordPosition records a user's location update.
func (s *Server) RecordPosition(ctx context.Context, req *pb.RecordPositionRequest) (*pb.RecordPositionResponse, error) {
	// get the id from the request metadata
	// if the id is not present, generate a new one
	slog.Info("RecordPosition request received", "user_id", req.UserId)
	uid, err := uuid.NewV7()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate uuid")
	}

	userPostion := domain.UserPosition{
		UserID:    req.UserId,
		Reference: uid.String(),
		Location: domain.GeoPoint{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		CreatedAt: req.Timestamp.AsTime(),
		PhoneMeta: domain.PhoneMeta{
			DeviceID:   req.Metadata.DeviceId,
			Brand:      req.Metadata.Brand,
			Model:      req.Metadata.Model,
			OS:         req.Metadata.Os,
			AppVersion: req.Metadata.AppVersion,
			Carrier:    req.Metadata.Carrier,
			Battery:    int(req.Metadata.Battery),
		},
	}

	err = s.services.UserPositionService.RecordPosition(ctx, userPostion)

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
