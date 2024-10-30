package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/aparnasukesh/inter-communication/movie_booking"
	"github.com/aparnasukesh/payment-svc/config"
	"github.com/razorpay/razorpay-go"
)

type service struct {
	repo          Repository
	bookingClient movie_booking.BookingServiceClient
	cfg           config.Config
}

type Service interface {
	ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error)
	GetTransactionStatus(ctx context.Context, transactionID int32) (*Transaction, error)
	PaymentSuccess(ctx context.Context, req PaymentStatusRequest) error
	PaymentFailure(ctx context.Context, req PaymentStatusRequest) error
}

func NewService(repo Repository, bookingClient movie_booking.BookingServiceClient, cfg config.Config) Service {
	return &service{
		repo:          repo,
		bookingClient: bookingClient,
		cfg:           cfg,
	}
}

func (s *service) GetTransactionStatus(ctx context.Context, transactionID int32) (*Transaction, error) {
	return s.repo.GetTransactionByID(ctx, transactionID)
}

func (s *service) ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error) {
	res, err := s.bookingClient.GetBookingByID(ctx, &movie_booking.GetBookingByIDRequest{
		BookingId: uint32(req.BookingID),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching booking: %v", err)
	}

	transaction := &Transaction{
		BookingID:       uint(req.BookingID),
		UserID:          uint(req.UserID),
		PaymentMethodID: 1,
		TransactionDate: time.Now(),
		Amount:          res.Booking.TotalAmount,
		Status:          "Pending",
	}

	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	client := razorpay.NewClient(s.cfg.RAZORPAY_KEY_ID, s.cfg.RAZORPAY_KEY_SECRET)
	data := map[string]interface{}{
		"amount":   int(res.Booking.TotalAmount * 100),
		"currency": "INR",
		"receipt":  fmt.Sprintf("txn_%d", transaction.TransactionID),
	}

	order, err := client.Order.Create(data, nil)
	if err != nil {
		transaction.Status = "Failed"
		if updateErr := s.repo.UpdateTransaction(ctx, transaction); updateErr != nil {
			return nil, fmt.Errorf("failed to update transaction status to failed: %v", updateErr)
		}
		_, _ = s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
			BookingId: int32(req.BookingID),
		})
		return nil, fmt.Errorf("error creating Razorpay order: %v", err)
	}

	fmt.Printf("Razorpay Order Created: %v\n", order)

	transaction.OrderID = order["id"].(string)
	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to update transaction with order ID: %v", err)
	}

	return transaction, nil
}

func (s *service) PaymentSuccess(ctx context.Context, req PaymentStatusRequest) error {
	transaction, err := s.repo.GetTransactionByOrderID(ctx, req.OrderID)
	if err != nil {
		return err
	}
	transaction.Status = "Success"
	transaction.RazorpayPaymentID = req.RazorpayPaymentID
	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
		return err
	}
	_, err = s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
		BookingId: int32(transaction.BookingID),
		Status:    "Success",
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) PaymentFailure(ctx context.Context, req PaymentStatusRequest) error {
	transaction, err := s.repo.GetTransactionByOrderID(ctx, req.OrderID)
	if err != nil {
		return err
	}
	transaction.Status = "Success"
	transaction.RazorpayPaymentID = req.RazorpayPaymentID
	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
		return err
	}
	_, err = s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
		BookingId: int32(transaction.BookingID),
		Status:    "Failed",
	})
	if err != nil {
		return err
	}
	return nil
}
