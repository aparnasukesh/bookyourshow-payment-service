package payment

import (
	"context"
	"encoding/json"
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
func (h *GrpcHandler) HandleRazorpayWebhook(ctx context.Context, req *payment.HandleRazorpayWebhookRequest) (*payment.HandleRazorpayWebhookResponse, error) {
	payload, err := json.Marshal(req.Payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling RazorpayPayload: %v", err)
	}

	err = h.svc.HandleRazorpayWebhook(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &payment.HandleRazorpayWebhookResponse{
		Message: "Webhook processed successfully",
		Status:  "success",
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