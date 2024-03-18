/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	common "github.com/nawaltni/api/gen/go/nawalt/common/v1"
	pb "github.com/nawaltni/api/gen/go/nawalt/tracker/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	addr   *string // address to connect to
	client *string // client id to use
	user   *string // user to use
	route  *int
)

// positionCmd represents the position command
var positionCmd = &cobra.Command{
	Use:   "position",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("position called")

		points := []string{}

		route1 := []string{
			"12.1012232, -86.2589056",
			"12.102133, -86.258303",
			"12.102842, -86.258323",
			"12.102774, -86.260305",
			"12.102538, -86.262668",
			"12.100684, -86.262536",
			"12.098991, -86.262870",
			"12.098044, -86.263346",
			"12.098928, -86.262864",
			"12.101572, -86.262579",
			"12.102484, -86.262297",
			"12.102731, -86.258305",
			"12.101400, -86.258246",
			"12.1002976, -86.2601484",
		}

		route2 := []string{
			"12.1011823, -86.2587094",
			"12.102721, -86.258322",
			"12.102715, -86.257023",
			"12.103047, -86.255228",
			"12.103902, -86.252854",
			"12.105042, -86.250872",
			"12.105987, -86.248909",
			"12.106538, -86.249088",
			"12.105398, -86.250470",
			"12.104453, -86.252854",
			"12.102821, -86.257812",
			"12.102716, -86.258305",
			"12.101762, -86.258284",
			"12.1019347, -86.2585626",
		}

		switch *route {
		case 1:
			points = route1
		case 2:
			points = route2
		}

		// route3 := []string{}

		uuid := uuid.New().String()

		// Set up a connection to the server.
		conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewTrackingServiceClient(conn)
		i := 0

		for {
			select {
			case <-time.After(5 * time.Second):
				// prepare request
				coords := strings.Split(points[i], ", ")
				fmt.Printf("%s - sending position %v\n", time.Now(), coords)
				lat, err := strconv.ParseFloat(coords[0], 64)
				if err != nil {
					log.Fatalf("could not parse latitude: %v", err)
				}
				lon, err := strconv.ParseFloat(coords[1], 64)
				if err != nil {
					log.Fatalf("could not parse longitude: %v", err)
				}

				req := &pb.RecordPositionRequest{
					RequestId: uuid,
					UserId:    *user,
					Timestamp: timestamppb.New(time.Now()),
					Location: &common.GeoPoint{
						Latitude:  float32(lat),
						Longitude: float32(lon),
					},
					// Timestamp: proto.TimestampNow(),
					ClientId: *client,
					Metadata: &pb.PhoneMetadata{
						DeviceId:   "123456",
						Brand:      "Samsung",
						Model:      "Galaxy S10",
						Os:         "Android",
						AppVersion: "0.0.1",
						Carrier:    "Claro",
						Battery:    100,
					},
				}
				// send request
				_, err = c.RecordPosition(context.Background(), req)
				if err != nil {
					log.Fatalf("could not position: %v", err)
				}
				i++
				if i == len(points) {
					i = 0
				}
			}
		}
	},
}

func init() {
	fakeCmd.AddCommand(positionCmd)

	addr = positionCmd.Flags().StringP("addr", "a", "localhost:50051", "Set the address to connect to")
	client = positionCmd.Flags().StringP("client", "c", "client1", "Set the client id to use")
	user = positionCmd.Flags().StringP("user", "u", "018c28a6-e56b-7f3f-a76f-76bc67373895", "Set the user to use")
	route = positionCmd.Flags().IntP("route", "r", 1, "Set the route to use")
	// make a request to the server every 60 seconds. We will use a predefined list of coordinates
}
