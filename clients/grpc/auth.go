package grpc

import (
	"context"
	"fmt"

	pb "github.com/nawaltni/api/gen/go/nawalt/auth/v1"
	"github.com/nawaltni/tracker/config"
	"github.com/nawaltni/tracker/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClientGRPC struct {
	config config.Auth
	client pb.AuthServiceClient
}

// NewAuthClientGRPC returns a new instance of AdapultClientGRPC
func NewAuthClientGRPC(conf config.Auth) (*AuthClientGRPC, error) {
	// create connection
	connString := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	conn, err := grpc.Dial(connString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// create client
	client := pb.NewAuthServiceClient(conn)

	return &AuthClientGRPC{
		config: conf,
		client: client,
	}, nil
}

// GetUserByBackendID calls the GetUserByBackendID method of the auth gRPC service
func (c *AuthClientGRPC) GetUserByBackendID(ctx context.Context, id int) (*domain.User, error) {
	// call the gRPC service
	req := &pb.GetUserByBackendIDRequest{
		BackendId: int32(id),
	}
	resp, err := c.client.GetUserByBackendID(ctx, req)
	if err != nil {
		return nil, err
	}

	res := MapUserFromProto(resp.User)

	return res, nil
}
