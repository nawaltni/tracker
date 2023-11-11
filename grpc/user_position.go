package grpc

import (
	"context"

	pb "github.com/nawaltni/api/gen/go/nawalt/tracker/v1"
	"github.com/pkg/errors"
)

// NewServer creates a new TrackingService server.
func NewServer() *Server {
	return &Server{}
}

// RecordPosition records a user's location update.
func (s *Server) RecordPosition(ctx context.Context, req *pb.RecordPositionRequest) (*pb.RecordPositionResponse, error) {
	// Implement your logic to handle the position recording.
	// You'll likely want to interact with a database/repository here.
	return nil, errors.New("not implemented")
}

// GetUserPosition retrieves a user's current position.
func (s *Server) GetUserPosition(ctx context.Context, req *pb.GetUserPositionRequest) (*pb.GetUserPositionResponse, error) {
	// Implement your logic to retrieve the user's current position.
	return nil, errors.New("not implemented")
}

// GetUserPositionList retrieves a user's position history.
func (s *Server) GetUserPositionList(ctx context.Context, req *pb.GetUserPositionListRequest) (*pb.GetUserPositionListResponse, error) {
	// Implement your logic to retrieve the user's position history.
	return nil, errors.New("not implemented")
}
