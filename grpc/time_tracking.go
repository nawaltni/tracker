package grpc

import (
	"context"
	"log/slog"

	pb "github.com/nawaltni/api/gen/go/nawalt/tracker/v1"
	"github.com/nawaltni/tracker/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CheckIn checks in a user to start a time tracking session
func (s *Server) CheckIn(ctx context.Context, req *pb.CheckInRequest) (*pb.CheckInResponse, error) {
	session := domain.TimeTrackingSession{
		UserID: req.UserId,
	}

	if req.Location != nil {
		session.LastKnownLocation = domain.GeoPoint{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		}
	}

	err := s.services.TimeTrackingService.CheckIn(ctx, &session)
	if err != nil {
		return &pb.CheckInResponse{}, status.Error(codes.Internal, "failed to check in")
	}

	return &pb.CheckInResponse{
		Status: convertSessionToProto(session),
	}, nil
}

// CheckOut checks out a user to end a time tracking session
func (s *Server) CheckOut(ctx context.Context, req *pb.CheckOutRequest) (*pb.CheckOutResponse, error) {
	session := domain.TimeTrackingSession{
		UserID:    req.UserId,
		SessionID: req.SessionId,
	}

	if req.Location != nil {
		session.LastKnownLocation = domain.GeoPoint{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		}
	}

	err := s.services.TimeTrackingService.CheckOut(ctx, &session)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check out")
	}

	return &pb.CheckOutResponse{
		Status: convertSessionToProto(session),
	}, nil
}

// StartBreak starts a break for the user within a time tracking session
func (s *Server) StartBreak(ctx context.Context, req *pb.StartBreakRequest) (*pb.StartBreakResponse, error) {
	session := domain.TimeTrackingSession{
		UserID:    req.UserId,
		SessionID: req.SessionId,
	}

	if req.Location != nil {
		session.LastKnownLocation = domain.GeoPoint{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		}
	}

	err := s.services.TimeTrackingService.StartBreak(ctx, &session)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start break")
	}

	return &pb.StartBreakResponse{
		Status: convertSessionToProto(session),
	}, nil
}

// EndBreak ends a break for the user within a time tracking session
func (s *Server) EndBreak(ctx context.Context, req *pb.EndBreakRequest) (*pb.EndBreakResponse, error) {
	session := domain.TimeTrackingSession{
		UserID:    req.UserId,
		SessionID: req.SessionId,
	}

	if req.Location != nil {
		session.LastKnownLocation = domain.GeoPoint{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		}
	}

	err := s.services.TimeTrackingService.EndBreak(ctx, &session)
	if err != nil {
		slog.Error("failed to end break", "error", err)
		return nil, status.Error(codes.Internal, "failed to end break")
	}

	return &pb.EndBreakResponse{
		Status: convertSessionToProto(session),
	}, nil
}

// GetTimeTrackingStatus gets the current time tracking status of the user
func (s *Server) GetTimeTrackingStatus(ctx context.Context, req *pb.GetTimeTrackingStatusRequest) (*pb.GetTimeTrackingStatusResponse, error) {
	session, err := s.services.TimeTrackingService.GetTimeTrackingSession(ctx, req.SessionId, req.UserId) // Assuming the session ID is not needed or can be retrieved differently
	if err != nil {
		return nil, status.Error(codes.NotFound, "session not found")
	}

	return &pb.GetTimeTrackingStatusResponse{
		Status: convertSessionToProto(*session),
	}, nil
}
