package payment

import (
	"context"
	"encoding/json"
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

// GetTransactionStatus retrieves the status of a payment transaction.

type Service interface {
	ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error)
	GetTransactionStatus(ctx context.Context, transactionID int32) (*Transaction, error)
	HandleRazorpayWebhook(ctx context.Context, payload []byte) error
}

func NewService(repo Repository, bookingClient movie_booking.BookingServiceClient, cfg config.Config) Service {
	return &service{
		repo:          repo,
		bookingClient: bookingClient,
		cfg:           cfg,
	}
}

// func (s *service) ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error) {
// 	res, err := s.bookingClient.GetBookingByID(ctx, &movie_booking.GetBookingByIDRequest{
// 		BookingId: uint32(req.BookingID),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	transaction := &Transaction{
// 		BookingID:       req.BookingID,
// 		UserID:          req.UserID,
// 		PaymentMethodID: req.PaymentMethodID,
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
// 		s.repo.UpdateTransaction(ctx, transaction)
// 		_, _ = s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
// 			BookingId: int32(req.BookingID),
// 		})
// 		return nil, fmt.Errorf("error creating Razorpay order: %v", err)
// 	}

// 	if order["status"] == "paid" {
// 		transaction.Status = "Success"
// 		s.repo.UpdateTransaction(ctx, transaction)
// 		_, err := s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
// 			BookingId: int32(req.BookingID),
// 			Status:    "Success",
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		transaction.Status = "Failed"
// 		s.repo.UpdateTransaction(ctx, transaction)
// 		_, err := s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
// 			BookingId: int32(req.BookingID),
// 			Status:    "Failed",
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 		_, _ = s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
// 			BookingId: int32(req.BookingID),
// 		})
// 	}

// 	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
// 		return nil, fmt.Errorf("failed to update transaction: %v", err)
// 	}
// 	return transaction, nil
// }

// func (s *service) ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error) {
// 	// Fetch booking details
// 	res, err := s.bookingClient.GetBookingByID(ctx, &movie_booking.GetBookingByIDRequest{
// 		BookingId: uint32(req.BookingID),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get booking details: %v", err)
// 	}

// 	// Create a new transaction with a pending status
// 	transaction := &Transaction{
// 		BookingID:       req.BookingID,
// 		UserID:          req.UserID,
// 		PaymentMethodID: req.PaymentMethodID,
// 		TransactionDate: time.Now(),
// 		Amount:          res.Booking.TotalAmount,
// 		Status:          "Pending",
// 	}

// 	// Save the transaction in the database
// 	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
// 		return nil, fmt.Errorf("failed to create transaction: %v", err)
// 	}

// 	// Initialize Razorpay client
// 	client := razorpay.NewClient(s.cfg.RAZORPAY_KEY_ID, s.cfg.RAZORPAY_KEY_SECRET)

// 	// Create a Razorpay order
// 	data := map[string]interface{}{
// 		"amount":   int(res.Booking.TotalAmount * 100), // Amount in paise (smallest currency unit in INR)
// 		"currency": "INR",
// 		"receipt":  fmt.Sprintf("txn_%d", transaction.TransactionID),
// 	}
// 	order, err := client.Order.Create(data, nil)
// 	if err != nil {
// 		// Update transaction as failed and delete the booking if order creation fails
// 		transaction.Status = "Failed"
// 		s.repo.UpdateTransaction(ctx, transaction)
// 		_, _ = s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
// 			BookingId: int32(req.BookingID),
// 		})
// 		return nil, fmt.Errorf("error creating Razorpay order: %v", err)
// 	}

// 	// Here, the status is still 'created', not 'paid'. Return the order ID to the user.
// 	// The actual payment status will be updated later via webhook or verification.
// 	transaction.Status = "Pending"
// 	transaction.OrderID = order["id"].(string) // Store the order ID for future reference

// 	// Update transaction in the database
// 	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
// 		return nil, fmt.Errorf("failed to update transaction: %v", err)
// 	}

//		// Return the transaction with the pending status and order ID
//		return transaction, nil
//	}
//
// Service method to process payment
// func (s *service) ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error) {
// 	// Step 1: Fetch booking details
// 	res, err := s.bookingClient.GetBookingByID(ctx, &movie_booking.GetBookingByIDRequest{
// 		BookingId: uint32(req.BookingID),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching booking: %v", err)
// 	}

// 	// Step 2: Create a transaction in local DB
// 	transaction := &Transaction{
// 		BookingID:       uint(req.BookingID),
// 		UserID:          uint(req.UserID),
// 		PaymentMethodID: 1, // Example payment method
// 		TransactionDate: time.Now(),
// 		Amount:          res.Booking.TotalAmount,
// 		Status:          "Pending", // Initially set as pending
// 	}

// 	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
// 		return nil, fmt.Errorf("failed to create transaction: %v", err)
// 	}

// 	// Step 3: Create a Razorpay order
// 	client := razorpay.NewClient(s.cfg.RAZORPAY_KEY_ID, s.cfg.RAZORPAY_KEY_SECRET)
// 	data := map[string]interface{}{
// 		"amount":   int(res.Booking.TotalAmount * 100), // Convert amount to paise (smallest unit)
// 		"currency": "INR",
// 		"receipt":  fmt.Sprintf("txn_%d", transaction.TransactionID),
// 	}

// 	order, err := client.Order.Create(data, nil)
// 	if err != nil {
// 		transaction.Status = "Failed"
// 		s.repo.UpdateTransaction(ctx, transaction)
// 		_, _ = s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
// 			BookingId: int32(req.BookingID),
// 		})
// 		return nil, fmt.Errorf("error creating Razorpay order: %v", err)
// 	}

// 	// Log order details
// 	fmt.Printf("Razorpay Order Created: %v\n", order)

// 	// Store Razorpay order ID in the transaction
// 	transaction.OrderID = order["id"].(string)
// 	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
// 		return nil, fmt.Errorf("failed to update transaction with order ID: %v", err)
// 	}

// 	// Step 4: Fetch payment status
// 	paymentStatus, err := client.Payment.Fetch(transaction.OrderID, map[string]interface{}{}, map[string]string{})
// 	if err != nil {
// 		transaction.Status = "Failed"
// 		s.repo.UpdateTransaction(ctx, transaction)
// 		return nil, fmt.Errorf("failed to fetch payment status: %v", err)
// 	}

// 	// Step 5: Handle payment success or failure
// 	if paymentStatus["status"] == "captured" { // Payment successful
// 		transaction.Status = "Success"
// 		if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
// 			return nil, fmt.Errorf("failed to update transaction status to success: %v", err)
// 		}
// 		// Update booking status as successful
// 		_, err := s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
// 			BookingId: int32(req.BookingID),
// 			Status:    "Success",
// 		})
// 		if err != nil {
// 			return nil, fmt.Errorf("error updating booking status: %v", err)
// 		}
// 	} else { // Payment failed
// 		transaction.Status = "Failed"
// 		s.repo.UpdateTransaction(ctx, transaction)
// 		_, err := s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
// 			BookingId: int32(req.BookingID),
// 			Status:    "Failed",
// 		})
// 		if err != nil {
// 			return nil, fmt.Errorf("error updating booking status: %v", err)
// 		}
// 		_, _ = s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
// 			BookingId: int32(req.BookingID),
// 		})
// 	}

// 	// Step 6: Return transaction details
// 	return transaction, nil
// }

func (s *service) GetTransactionStatus(ctx context.Context, transactionID int32) (*Transaction, error) {
	return s.repo.GetTransactionByID(ctx, transactionID)
}
func (s *service) ProcessPayment(ctx context.Context, req *PaymentRequest) (*Transaction, error) {
	// Step 1: Fetch booking details
	res, err := s.bookingClient.GetBookingByID(ctx, &movie_booking.GetBookingByIDRequest{
		BookingId: uint32(req.BookingID),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching booking: %v", err)
	}

	// Step 2: Create a transaction in the local DB
	transaction := &Transaction{
		BookingID:       uint(req.BookingID),
		UserID:          uint(req.UserID),
		PaymentMethodID: 1, // Example payment method
		TransactionDate: time.Now(),
		Amount:          res.Booking.TotalAmount,
		Status:          "Pending", // Initially set as pending
	}

	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Step 3: Create a Razorpay order
	client := razorpay.NewClient(s.cfg.RAZORPAY_KEY_ID, s.cfg.RAZORPAY_KEY_SECRET)
	data := map[string]interface{}{
		"amount":   int(res.Booking.TotalAmount * 100), // Convert amount to paise (smallest unit)
		"currency": "INR",
		"receipt":  fmt.Sprintf("txn_%d", transaction.TransactionID),
	}

	order, err := client.Order.Create(data, nil)
	if err != nil {
		transaction.Status = "Failed"
		if updateErr := s.repo.UpdateTransaction(ctx, transaction); updateErr != nil {
			return nil, fmt.Errorf("failed to update transaction status to failed: %v", updateErr)
		}
		// Clean up the booking if the order creation fails
		_, _ = s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
			BookingId: int32(req.BookingID),
		})
		return nil, fmt.Errorf("error creating Razorpay order: %v", err)
	}

	// Log order details
	fmt.Printf("Razorpay Order Created: %v\n", order)

	// Store Razorpay order ID in the transaction
	transaction.OrderID = order["id"].(string)
	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to update transaction with order ID: %v", err)
	}

	// Step 4: Return transaction details
	return transaction, nil
}
func (s *service) HandleRazorpayWebhook(ctx context.Context, payload []byte) error {
	var webhookEvent RazorpayWebhookPayload

	// Unmarshal the webhook payload
	if err := json.Unmarshal(payload, &webhookEvent); err != nil {
		return fmt.Errorf("error unmarshalling webhook payload: %v", err)
	}

	// Fetch the transaction using Order ID
	transaction, err := s.repo.GetTransactionByOrderID(ctx, webhookEvent.Data.OrderID)
	if err != nil {
		return fmt.Errorf("transaction not found for order ID %s: %v", webhookEvent.Data.OrderID, err)
	}

	// Update transaction details
	transaction.Status = webhookEvent.Data.Status
	transaction.Amount = webhookEvent.Data.Amount

	// Handle the payment status manually
	switch webhookEvent.Data.Status {
	case "Pending":
		// Simulate payment confirmation logic (you could call an external service or just simulate here)
		isPaymentSuccessful := simulatePaymentConfirmation(transaction)

		if isPaymentSuccessful {
			transaction.Status = "Success"
			if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
				return fmt.Errorf("failed to update transaction status to success: %v", err)
			}
			// Update booking status as successful
			_, err = s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
				BookingId: int32(transaction.BookingID),
				Status:    "Success",
			})
			if err != nil {
				return fmt.Errorf("error updating booking status: %v", err)
			}
		} else {
			transaction.Status = "Failed"
			if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
				return fmt.Errorf("failed to update transaction status to failed: %v", err)
			}
			// Update booking status as failed
			_, err = s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
				BookingId: int32(transaction.BookingID),
				Status:    "Failed",
			})
			if err != nil {
				return fmt.Errorf("error updating booking status: %v", err)
			}
			_, err := s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
				BookingId: int32(transaction.BookingID),
			})
			if err != nil {
				return err
			}
		}

	case "Success":
		transaction.Status = "Success"
		if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
			return fmt.Errorf("failed to update transaction status to success: %v", err)
		}
		// Update booking status as successful
		_, err = s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
			BookingId: int32(transaction.BookingID),
			Status:    "Success",
		})
		if err != nil {
			return fmt.Errorf("error updating booking status: %v", err)
		}

	case "Failed":
		transaction.Status = "Failed"
		if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
			return fmt.Errorf("failed to update transaction status to failed: %v", err)
		}
		// Update booking status as failed
		_, err = s.bookingClient.UpdateBookingStatusByBookingID(ctx, &movie_booking.UpdateBookingStatusByBookingIDRequest{
			BookingId: int32(transaction.BookingID),
			Status:    "Failed",
		})
		if err != nil {
			return fmt.Errorf("error updating booking status: %v", err)
		}
		_, err := s.bookingClient.DeleteBookingByBookingID(ctx, &movie_booking.DeleteBookingByIDRequest{
			BookingId: int32(transaction.BookingID),
		})
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unhandled payment status: %v", webhookEvent.Data.Status)
	}

	return nil
}

// Simulate a manual confirmation of the payment, this can be business logic or an API call
func simulatePaymentConfirmation(transaction *Transaction) bool {
	// Example simulation: Assume every payment with amount > 500 is successful for testing purposes
	if transaction.Amount < 500 {
		return true
	}
	return false
}
