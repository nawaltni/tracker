package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	pb "github.com/nawaltni/api/gen/go/nawalt/tracker/v1"
	"github.com/nawaltni/tracker/config"
	"github.com/nawaltni/tracker/services"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
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

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
			// Add any other interceptor you want.
		),
	)
	reflection.Register(s)
	grpc_health_v1.RegisterHealthServer(s, &server)
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

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// Check implements the health check for the GRPC server
func (s *Server) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

// Watch implements the health check for the GRPC server
func (s *Server) Watch(in *grpc_health_v1.HealthCheckRequest, _ grpc_health_v1.Health_WatchServer) error {
	// Example of how to register both methods but only implement the Check method.
	return nil
}
