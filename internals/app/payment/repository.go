package payment

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
	GetTransactionByID(ctx context.Context, transactionID int32) (*Transaction, error)
	GetTransactionByOrderID(ctx context.Context, orderId string) (*Transaction, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateTransaction(ctx context.Context, transaction *Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *repository) UpdateTransaction(ctx context.Context, transaction *Transaction) error {
	return r.db.WithContext(ctx).Model(&Transaction{}).Where("transaction_id = ?", transaction.TransactionID).Updates(transaction).Error
}

func (r *repository) GetTransactionByID(ctx context.Context, transactionID int32) (*Transaction, error) {
	var transaction Transaction
	err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&transaction).Error
	return &transaction, err
}

func (r *repository) GetTransactionByOrderID(ctx context.Context, orderId string) (*Transaction, error) {
	transaction := &Transaction{}
	if err := r.db.Where("order_id = ?", orderId).First(&transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}
