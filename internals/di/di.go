package di

import (
	"log"

	"github.com/aparnasukesh/payment-svc/config"
	"github.com/aparnasukesh/payment-svc/internals/app/payment"
	"github.com/aparnasukesh/payment-svc/internals/boot"
	grpcclient "github.com/aparnasukesh/payment-svc/pkg/grpcClient"
	"github.com/aparnasukesh/payment-svc/pkg/sql"
)

func InitResources(cfg config.Config) (func() error, error) {
	// Db initialization
	db, err := sql.NewSql(cfg)
	if err != nil {
		log.Fatal(err)
	}
	_, _, movieBookingClient, err := grpcclient.NewMovieBookingGrpcClint(cfg.GrpcMovieBookingServicePort)
	if err != nil {
		return nil, err
	}
	repo := payment.NewRepository(db)
	service := payment.NewService(repo, movieBookingClient)
	handler := payment.NewGrpcHandler(service)

	server, err := boot.NewGrpcServer(cfg, handler)
	if err != nil {
		log.Fatal(err)
	}
	return server, nil
}
