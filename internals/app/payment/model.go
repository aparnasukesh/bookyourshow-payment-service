package payment

import (
	"time"
)

type PaymentMethod struct {
	PaymentMethodID uint   `gorm:"primaryKey;autoIncrement" json:"payment_method_id"`
	MethodName      string `gorm:"type:varchar(100);not null" json:"method_name"`
}
type Transaction struct {
	TransactionID   uint      `gorm:"primaryKey;autoIncrement" json:"transaction_id"`
	BookingID       uint      `gorm:"not null" json:"booking_id"`
	UserID          uint      `gorm:"not null" json:"user_id"`
	PaymentMethodID uint      `gorm:"not null" json:"payment_method_id"`
	TransactionDate time.Time `gorm:"type:timestamp;not null" json:"transaction_date"`
	Amount          float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status          string    `gorm:"type:varchar(50);not null" json:"status"`
	OrderID         string    `gorm:"type:varchar(100)" json:"order_id"` // Add this field
}

type PaymentRequest struct {
	BookingID       uint      `gorm:"not null" json:"booking_id"`
	UserID          uint      `gorm:"not null" json:"user_id"`
	PaymentMethodID uint      `gorm:"not null" json:"payment_method_id"`
	TransactionDate time.Time `gorm:"type:timestamp;not null" json:"transaction_date"`
	Amount          float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
}

type RazorpayWebhookPayload struct {
	Data struct {
		TransactionID   int     `json:"transaction_id"`
		BookingID       int     `json:"booking_id"`
		UserID          int     `json:"user_id"`
		PaymentMethodID int     `json:"payment_method_id"`
		TransactionDate string  `json:"transaction_date"` // You may want to use time.Time for better handling
		Amount          float64 `json:"amount"`
		OrderID         string  `json:"order_id"`
		Status          string  `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}
