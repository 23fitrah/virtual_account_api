package validations

type CreateVirtualAccountValidation struct {
	UserValidation
	CustomerID   string  `json:"customer_id" validate:"required"`
	CustomerName string  `json:"customer_name" validate:"required"`
	Amount       float64 `json:"amount" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	ReferenceID  string  `json:"reference_id" validate:"required"`
}

type DataVirtualAccountValidation struct {
	CreateVirtualAccountValidation
	IdVa           string  `json:"idVa" validate:"required"`
	MerchantNumber string  `json:"merchantNumber" validate:"required"`
	CustReference  string  `json:"customerReference" validate:"required"`
	Amount         float64 `json:"amount" validate:"required"`
	Status         string  `json:"status" validate:"required"`
	ExpiredAt      string  `json:"expiredAt" validate:"required"`
}
