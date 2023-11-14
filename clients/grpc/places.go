package grpc

import (
	"context"
	"fmt"

	pb "github.com/nawaltni/api/gen/go/nawalt/places/v1"
	"github.com/nawaltni/tracker/config"
	"github.com/nawaltni/tracker/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PlacesClientGRPC struct {
	config config.Places
	client pb.PlacesServiceClient
}

// NewPlacesClientGRPC returns a new instance of AdapultClientGRPC
func NewPlacesClientGRPC(conf config.Places) (*PlacesClientGRPC, error) {
	// create connection
	connString := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	conn, err := grpc.Dial(connString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// create client
	client := pb.NewPlacesServiceClient(conn)

	return &PlacesClientGRPC{
		config: conf,
		client: client,
	}, nil
}

// IsWithinPlace checks if the given coordinates are within a place by calling the places
// GRPC service.
func (c *PlacesClientGRPC) IsWithinPlace(ctx context.Context, lat float32, lon float32) (*domain.WithinPlaceResponse, error) {
	// call the gRPC service
	req := &pb.IsWithinPlaceRequest{
		Latitude:  lat,
		Longitude: lon,
	}
	resp, err := c.client.IsWithinPlace(ctx, req)
	if err != nil {
		return nil, err
	}

	res := domain.WithinPlaceResponse{
		IsWithin:  resp.IsWithin,
		PlaceID:   resp.PlaceId,
		PlaceName: resp.Name,
	}

	return &res, nil
}
