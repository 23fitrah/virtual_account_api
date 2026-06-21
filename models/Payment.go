package models

import "time"

type Payment struct {
	ID               string    `gorm:"column:ID;primaryKey"`
	VANumber         string    `gorm:"column:VA_NUMBER"`
	VirtualAccountID string    `gorm:"column:VIRTUAL_ACCOUNT_ID"`
	PaidAmount       float64   `gorm:"column:PAID_AMOUNT"`
	PaymentChannel   string    `gorm:"column:PAYMENT_CHANNEL"`
	Status           int       `gorm:"column:STATUS"`
	RawPayload       string    `gorm:"column:RAW_PAYLOAD"`
	CreatedAt        time.Time `gorm:"column:CREATED_AT"`
}

func (Payment) TableName() string {
	return "payments"
}
