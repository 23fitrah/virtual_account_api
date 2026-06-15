package resources

import (
	"time"
	"virtual_account_api/models"
)

type GetVAResource struct {
	*DataGetVAStatus
}
type DataGetVAStatus struct {
	Id        string `json:"id"`
	VANumber  string `json:"va_number"`
	Status    string `json:"status"`
	Amount    string `json:"amount"`
	ExpiredAt string `json:"expired_at"`
}

type CreateVAResource struct {
	ID           string    `json:"id"`
	VANumber     string    `json:"va_number"`
	CustomerID   string    `json:"customer_id"`
	CustomerName string    `json:"customer_name"`
	Amount       float64   `json:"amount"`
	Description  string    `json:"description"`
	ReferenceID  string    `json:"reference_id"`
	ExpiredAt    time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
}

type DataVAStatus struct {
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

type ListVAResource struct {
	Data []GetVAListResource `json:"data"`
}

type GetVAListResource struct {
	Id           string `json:"id"`
	VANumber     string `json:"va_number"`
	CustomerId   string `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	Amount       string `json:"amount"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	ReferenceId  string `json:"reference_id"`
	ExpiredAt    string `json:"expired_at"`
	PaidAt       string `json:"paid_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func ToFormModelResource(data *models.VirtualAccount) CreateVAResource {
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
	return CreateVAResource{
		ID:           data.ID,
		VANumber:     data.VANumber,
		CustomerID:   data.CustomerID,
		CustomerName: data.CustomerName,
		Amount:       data.Amount,
		Description:  data.Description,
		ReferenceID:  data.ReferenceID,
		ExpiredAt:    data.ExpiredAt,
		CreatedAt:    data.CreatedAt,
		Status:       status,
	}
}
