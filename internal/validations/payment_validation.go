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

/*
	type DataVirtualAccountValidation struct {
		CreateVAValidation
		IdVa           string  `json:"idVa" validate:"required"`
		MerchantNumber string  `json:"merchantNumber" validate:"required"`
		CustReference  string  `json:"customerReference" validate:"required"`
		Amount         float64 `json:"amount" validate:"required"`
		Status         string  `json:"status" validate:"required"`
		ExpiredAt      string  `json:"expiredAt" validate:"required"`
	}
*/
type GetPaymentHistoryValidation struct {
	UserValidation
}

/*
type GetVAStatusValidation struct {
	UserValidation
}*/
