package payment

import (
	"context"
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
	res, err := h.svc.ProcessPayment(ctx, &PaymentRequest{
		BookingID:       uint(req.BookingId),
		UserID:          uint(req.UserId),
		PaymentMethodID: uint(req.PaymentMethodId),
		TransactionDate: time.Time{},
		Amount:          req.Amount,
	})
	if err != nil {
		return nil, err
	}
	return &payment.ProcessPaymentResponse{
		Transaction: &payment.Transaction{
			TransactionID:   int32(res.TransactionID),
			BookingID:       req.BookingId,
			UserID:          int32(res.UserID),
			PaymentMethodID: req.PaymentMethodId,
			TransactionDate: res.TransactionDate.String(),
			Amount:          req.Amount,
			Status:          res.Status,
		},
	}, nil
}
