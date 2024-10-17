package payment

import (
	"context"
	"time"

	"github.com/aparnasukesh/inter-communication/movie_booking"
)

type service struct {
	repo          Repository
	bookingClient movie_booking.BookingServiceClient
}

type Service interface {
	ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error)
}

func NewService(repo Repository, bookingClient movie_booking.BookingServiceClient) Service {
	return &service{
		repo:          repo,
		bookingClient: bookingClient,
	}
}

func (s *service) ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error) {
	res, err := s.bookingClient.GetBookingByID(ctx, &movie_booking.GetBookingByIDRequest{
		BookingId: uint32(req.BookingID),
	})
	if err != nil {
		return nil, err
	}
	s.repo
	return &Transaction{
		TransactionID:   res.,
		BookingID:       0,
		UserID:          0,
		PaymentMethodID: 0,
		TransactionDate: time.Time{},
		Amount:          0,
		Status:          "",
	}
}
