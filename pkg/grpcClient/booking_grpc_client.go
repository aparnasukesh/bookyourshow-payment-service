package grpcclient

import (
	"log"

	pb "github.com/aparnasukesh/inter-communication/movie_booking"
	"google.golang.org/grpc"
)

func NewMovieBookingGrpcClint(port string) (pb.MovieServiceClient, pb.TheatreServiceClient, pb.BookingServiceClient, error) {
	// conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	// if err != nil {
	// 	return nil, nil, nil, err
	// }
	address := "movies-booking-svc.default.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, nil, nil, err
	}
	return pb.NewMovieServiceClient(conn), pb.NewTheatreServiceClient(conn), pb.NewBookingServiceClient(conn), nil
}
