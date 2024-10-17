package grpcclient

import (
	pb "github.com/aparnasukesh/inter-communication/movie_booking"
	"google.golang.org/grpc"
)

func NewMovieBookingGrpcClint(port string) (pb.MovieServiceClient, pb.TheatreServiceClient, pb.BookingServiceClient, error) {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, nil, nil, err
	}
	return pb.NewMovieServiceClient(conn), pb.NewTheatreServiceClient(conn), pb.NewBookingServiceClient(conn), nil
}
