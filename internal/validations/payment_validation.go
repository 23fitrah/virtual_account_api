package validations

type CallbackPaymentValidation struct {
	UserValidation
	RequestData FormPaymentValidation `json:"payload" validate:"required"`
}

type FormPaymentValidation struct {
	VANumber       string  `json:"va_number" validate:"required"`
	PaidAmount     float64 `json:"paid_amount" validate:"required"`
	PaymentChannel string  `json:"payment_channel" validate:"required"`
	ReferenceID    string  `json:"reference_id" validate:"required"`
}

type GetPaymentHistoryValidation struct {
	UserValidation
}
