package grpc

import (
	"fmt"
	"log/slog"
	"net"

	pb "github.com/nawaltni/api/gen/go/nawalt/tracker/v1"
	"github.com/nawaltni/tracker/config"
	"github.com/nawaltni/tracker/services"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is the gRPC server
type Server struct {
	pb.UnimplementedTrackingServiceServer
	conn     net.Listener
	s        *grpc.Server
	services *services.Services
}

// New creates a new gRPC server for PlacesService
func New(conf config.Configuration, services *services.Services) (*Server, error) {
	serverHost := fmt.Sprintf("%s:%d", conf.GRPC.Host, conf.GRPC.Port)
	conn, err := net.Listen("tcp", serverHost)
	if err != nil {
		return nil, errors.Wrap(err, "failed to listen")
	}

	server := Server{
		conn:     conn,
		services: services,
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterTrackingServiceServer(s, &server)
	server.s = s

	return &server, nil
}

// Start starts the gRPC server
func (s *Server) Start() error {
	slog.Info("Starting Places gRPC server")
	if err := s.s.Serve(s.conn); err != nil {
		return errors.Wrap(err, "failed to start Places gRPC server")
	}

	return nil
}
