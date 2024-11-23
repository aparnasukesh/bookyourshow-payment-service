package boot

import (
	"log"
	"net"

	pay "github.com/aparnasukesh/inter-communication/payment"
	"github.com/aparnasukesh/payment-svc/config"
	"github.com/aparnasukesh/payment-svc/internals/app/payment"
	"google.golang.org/grpc"
)

func NewGrpcServer(config config.Config, paymentGrpcHandler payment.GrpcHandler) (func() error, error) {
	//lis, err := net.Listen("tcp", ":"+config.GrpcPort)
	lis, err := net.Listen("tcp", "0.0.0.0:"+config.GrpcPort)
	if err != nil {
		return nil, err
	}
	s := grpc.NewServer()
	pay.RegisterPaymentServiceServer(s, &paymentGrpcHandler)
	srv := func() error {
		log.Printf("gRPC server started on port %s", config.GrpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return err
		}
		return nil
	}
	return srv, nil
}
