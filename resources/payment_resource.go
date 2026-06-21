package resources

import (
	"time"
	"virtual_account_api/models"
)

type PaymentResource struct {
	Id         string `json:"id"`
	VANumber   string `json:"va_number"`
	StatusName string `json:"status_name"`
	Amount     string `json:"amount"`
	ExpiredAt  string `json:"expired_at"`
}

type CreatePaymentResource struct {
	VANumber string `json:"va_number"`
	Status   string `json:"status"`
}

type DataPaymentStatus struct {
	Id           string `json:"id"`
	VANumber     string `json:"va_number"`
	CustomerId   string `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	Amount       string `json:"amount"`
	Description  string `json:"description"`
	ReferenceId  string `json:"reference_id"`
	ExpiredAt    string `json:"expired_at"`
	CreatedAt    string `json:"created_at"`
}

type ListPaymentResource struct {
	Data []GetPaymentListResource `json:"payload"`
}

type GetPaymentListResource struct {
	Id             string    `json:"id"`
	VANumber       string    `json:"va_number"`
	PaidAmount     string    `json:"paid_amount"`
	PaymentChannel string    `json:"payment_channel"`
	StatusName     string    `json:"status_name"`
	CreatedAt      time.Time `json:"created_at"`
}

func ToFormModelPaymentResource(data *models.Payment) CreatePaymentResource {
	status := ""
	switch data.Status {
	case 1:
		status = "PENDING"
	case 2:
		status = "PAID"
	case 3:
		status = "EXPIRED"
	case 4:
		status = "CANCELED"
	}
	return CreatePaymentResource{
		VANumber: data.VANumber,
		Status:   status,
	}
}
