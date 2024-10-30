package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/aparnasukesh/inter-communication/payment"
)

type GrpcHandler struct {
	svc Service
	payment.UnimplementedPaymentServiceServer
}

func NewGrpcHandler(svc Service) GrpcHandler {
	return GrpcHandler{
		svc: svc,
	}
}
func (h *GrpcHandler) ProcessPayment(ctx context.Context, req *payment.ProcessPaymentRequest) (*payment.ProcessPaymentResponse, error) {
	paymentReq := &PaymentRequest{
		BookingID:       uint(req.BookingId),
		UserID:          uint(req.UserId),
		PaymentMethodID: uint(req.PaymentMethodId),
		TransactionDate: time.Now(),
		Amount:          req.Amount,
	}

	transaction, err := h.svc.ProcessPayment(ctx, paymentReq)
	if err != nil {
		return nil, fmt.Errorf("error processing payment: %v", err)
	}

	return &payment.ProcessPaymentResponse{
		Transaction: &payment.Transaction{
			TransactionId:   int32(transaction.TransactionID),
			BookingId:       int32(transaction.BookingID),
			UserId:          int32(transaction.UserID),
			PaymentMethodId: int32(transaction.PaymentMethodID),
			TransactionDate: transaction.TransactionDate.Format(time.RFC3339),
			Amount:          transaction.Amount,
			OrderId:         transaction.OrderID,
			Status:          transaction.Status,
		},
	}, nil
}

func (h *GrpcHandler) GetTransactionStatus(ctx context.Context, req *payment.GetTransactionStatusRequest) (*payment.GetTransactionStatusResponse, error) {
	transaction, err := h.svc.GetTransactionStatus(ctx, req.TransactionId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transaction status: %v", err)
	}
	return &payment.GetTransactionStatusResponse{
		TransactionId:   int32(transaction.TransactionID),
		Status:          transaction.Status,
		Amount:          transaction.Amount,
		PaymentMethodId: int32(transaction.PaymentMethodID),
		TransactionDate: transaction.TransactionDate.Format(time.RFC3339),
	}, nil
}

func (h *GrpcHandler) PaymentSuccess(ctx context.Context, req *payment.PaymentSuccessRequest) (*payment.PaymentSuccessResponse, error) {
	err := h.svc.PaymentSuccess(ctx, PaymentStatusRequest{
		BookingID:         int(req.BookingId),
		OrderID:           req.OrderId,
		RazorpayPaymentID: req.RazorpayPaymentId,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *GrpcHandler) PaymentFailure(ctx context.Context, req *payment.PaymentFailureRequest) (*payment.PaymentFailureResponse, error) {
	err := h.svc.PaymentFailure(ctx, PaymentStatusRequest{
		BookingID:         int(req.BookingId),
		OrderID:           req.OrderId,
		RazorpayPaymentID: req.RazorpayPaymentId,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
