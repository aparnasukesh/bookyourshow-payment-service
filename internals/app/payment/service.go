// package payment

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/aparnasukesh/inter-communication/movie_booking"
// 	"github.com/aparnasukesh/payment-svc/config"
// 	"github.com/razorpay/razorpay-go"
// )

// type service struct {
// 	repo          Repository
// 	bookingClient movie_booking.BookingServiceClient
// 	cfg           config.Config
// }

// type Service interface {
// 	ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error)
// 	GetTransactionStatus(ctx context.Context, transactionID int32) (*Transaction, error)
// 	PaymentSuccess(ctx context.Context, req PaymentStatusRequest) error
// 	PaymentFailure(ctx context.Context, req PaymentStatusRequest) error
// }

// func NewService(repo Repository, bookingClient movie_booking.BookingServiceClient, cfg config.Config) Service {
// 	return &service{
// 		repo:          repo,
// 		bookingClient: bookingClient,
// 		cfg:           cfg,
// 	}
// }

// func (s *service) GetTransactionStatus(ctx context.Context, transactionID int32) (*Transaction, error) {
// 	return s.repo.GetTransactionByID(ctx, transactionID)
// }

// func (s *service) ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error) {
// 	res, err := s.bookingClient.GetBookingByID(ctx, &movie_booking.GetBookingByIDRequest{
// 		BookingId: uint32(req.BookingID),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching booking: %v", err)
// 	}

// 	transaction := &Transaction{
// 		BookingID:       uint(req.BookingID),
// 		UserID:          uint(req.UserID),
// 		PaymentMethodID: 1,
// 		TransactionDate: time.Now(),
// 		Amount:          res.Booking.TotalAmount,
// 		Status:          "Pending",
// 	}

// 	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
// 		return nil, fmt.Errorf("failed to create transaction: %v", err)
// 	}

// 	client := razorpay.NewClient(s.cfg.RAZORPAY_KEY_ID, s.cfg.RAZORPAY_KEY_SECRET)
// 	data := map[string]interface{}{
// 		"amount":   int(res.Booking.TotalAmount * 100),
// 		"currency": "INR",
// 		"receipt":  fmt.Sprintf("txn_%d", transaction.TransactionID),
// 	}

// 	order, err := client.Order.Create(data, nil)
// 	if err != nil {
// 		transaction.Status = "Failed"
// 		if updateErr := s.repo.UpdateTransaction(ctx, transaction); updateErr != nil {
// 			return nil, fmt.Errorf("failed to update transaction status to failed: %v", updateErr)
// 		}
// 		_, _ = s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
// 			BookingId: int32(req.BookingID),
// 		})
// 		return nil, fmt.Errorf("error creating Razorpay order: %v", err)
// 	}

// 	fmt.Printf("Razorpay Order Created: %v\n", order)

// 	transaction.OrderID = order["id"].(string)
// 	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
// 		return nil, fmt.Errorf("failed to update transaction with order ID: %v", err)
// 	}

// 	return transaction, nil
// }

// func (s *service) PaymentSuccess(ctx context.Context, req PaymentStatusRequest) error {
// 	transaction, err := s.repo.GetTransactionByOrderID(ctx, req.OrderID)
// 	if err != nil {
// 		return err
// 	}
// 	transaction.Status = "Success"
// 	transaction.RazorpayPaymentID = req.RazorpayPaymentID
// 	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
// 		return err
// 	}
// 	_, err = s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
// 		BookingId: int32(transaction.BookingID),
// 		Status:    "Success",
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *service) PaymentFailure(ctx context.Context, req PaymentStatusRequest) error {
// 	transaction, err := s.repo.GetTransactionByOrderID(ctx, req.OrderID)
// 	if err != nil {
// 		return err
// 	}
// 	transaction.Status = "Success"
// 	transaction.RazorpayPaymentID = req.RazorpayPaymentID
// 	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
// 		return err
// 	}
// 	_, err = s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
// 		BookingId: int32(transaction.BookingID),
// 		Status:    "Failed",
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

package payment

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aparnasukesh/inter-communication/movie_booking"
	"github.com/aparnasukesh/payment-svc/config"
	"github.com/razorpay/razorpay-go"
)

// CircuitBreaker represents a simple circuit breaker implementation.
type CircuitBreaker struct {
	state          string
	failureCount   int
	successCount   int
	timeout        time.Duration
	errorThreshold int
	retryDelay     time.Duration
	lastFailure    time.Time
}

// Constants for circuit breaker states
const (
	StateClosed   = "CLOSED"
	StateOpen     = "OPEN"
	StateHalfOpen = "HALF_OPEN"
)

// NewCircuitBreaker initializes a new circuit breaker.
func NewCircuitBreaker(timeout time.Duration, errorThreshold int, retryDelay time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:          StateClosed,
		timeout:        timeout,
		errorThreshold: errorThreshold,
		retryDelay:     retryDelay,
	}
}

// Call executes a function within the context of the circuit breaker.
func (cb *CircuitBreaker) Call(fn func() (interface{}, error)) (interface{}, error) {
	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = StateHalfOpen
		} else {
			return nil, fmt.Errorf("circuit breaker is open")
		}
	case StateHalfOpen:
		// Allow a single request to test if the service is back
		cb.state = StateClosed
	default:
	}

	result, err := fn()
	if err != nil {
		cb.failureCount++
		cb.lastFailure = time.Now()
		if cb.failureCount >= cb.errorThreshold {
			cb.state = StateOpen
		}
		return nil, err
	}

	cb.successCount++
	cb.failureCount = 0 // Reset on success
	return result, nil
}

type service struct {
	repo          Repository
	bookingClient movie_booking.BookingServiceClient
	cfg           config.Config
	cbBooking     *CircuitBreaker
	cbPayment     *CircuitBreaker
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
		cbBooking:     NewCircuitBreaker(5*time.Second, 3, 2*time.Second), // Adjust these values as needed
		cbPayment:     NewCircuitBreaker(5*time.Second, 3, 2*time.Second),
	}
}

func (s *service) GetTransactionStatus(ctx context.Context, transactionID int32) (*Transaction, error) {
	return s.repo.GetTransactionByID(ctx, transactionID)
}

func (s *service) ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error) {
	var res *movie_booking.GetBookingByIDResponse

	// Attempt to call the booking service with a circuit breaker
	bookingResult, err := s.cbBooking.Call(func() (interface{}, error) {
		return s.bookingClient.GetBookingByID(ctx, &movie_booking.GetBookingByIDRequest{
			BookingId: uint32(req.BookingID),
		})
	})

	if err != nil {
		log.Printf("Error fetching booking: %v", err)
		return nil, fmt.Errorf("error fetching booking: %v", err)
	}
	res = bookingResult.(*movie_booking.GetBookingByIDResponse)

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
	for i := 0; i < 3; i++ {
		orderResult, err := s.cbPayment.Call(func() (interface{}, error) {
			return client.Order.Create(data, nil)
		})

		if err == nil {
			// Order creation was successful, so access the "id" field in the map
			if order, ok := orderResult.(map[string]interface{}); ok {
				if orderID, exists := order["id"].(string); exists {
					transaction.OrderID = orderID
					break
				} else {
					log.Println("Order ID not found in the Razorpay order response")
				}
			}
		} else {
			log.Printf("Retrying payment order creation: attempt %d", i+1)
			time.Sleep(s.cbPayment.retryDelay)
		}
	}

	if transaction.OrderID == "" {
		transaction.Status = "Failed"
		if updateErr := s.repo.UpdateTransaction(ctx, transaction); updateErr != nil {
			return nil, fmt.Errorf("failed to update transaction status to failed: %v", updateErr)
		}
		_, _ = s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
			BookingId: int32(req.BookingID),
		})
		return nil, fmt.Errorf("error creating Razorpay order: transaction failed")
	}

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
	transaction.Status = "Failed"
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
