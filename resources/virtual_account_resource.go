package resources

type GetVAResource struct {
	BaseResponse
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
	BaseResponse
	*DataCreateVAStatus
}
type DataCreateVAStatus struct {
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
	BaseResponse
	Data []GetVAListResource `json:"data"`
}

type GetVAListResource struct {
	Id           string `json:"id"`
	VANumber     string `json:"va_number"`
	CustomerId   string `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	Amount       string `json:"amount"`
	Description  string `json:"description"`
}
