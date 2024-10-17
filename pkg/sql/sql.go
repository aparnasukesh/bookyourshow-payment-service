package sql

import (
	"fmt"
	"log"
	"sync"

	"github.com/aparnasukesh/payment-svc/config"
	"github.com/aparnasukesh/payment-svc/internals/app/payment"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	mutex      sync.Mutex
	isExist    map[string]bool
)

func NewSql(config config.Config) (*gorm.DB, error) {
	if dbInstance == nil && !isExist[config.DBName] {
		mutex.Lock()
		defer mutex.Unlock()
		if dbInstance == nil && !isExist[config.DBName] {
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatal(err.Error())
				return nil, err
			}
			dbInstance = db
		}
	}
	err := dbInstance.AutoMigrate(&payment.PaymentMethod{})
	if err != nil {
		log.Fatalf("Error migrating PaymentMethod table: %v", err)
	}

	var count int64
	dbInstance.Model(&payment.PaymentMethod{}).Where("payment_method_id = ?", 1).Count(&count)

	if count == 0 {
		razorpay := payment.PaymentMethod{
			PaymentMethodID: 1,
			MethodName:      "Razorpay",
		}
		if err := dbInstance.Create(&razorpay).Error; err != nil {
			log.Fatalf("Error inserting Razorpay method: %v", err)
		} else {
			log.Println("Razorpay method inserted successfully.")
		}
	} else {
		log.Println("Razorpay method already exists. No need to insert.")
	}
	dbInstance.AutoMigrate(&payment.Transaction{})
	return dbInstance, nil
}
