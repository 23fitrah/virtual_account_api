package resources

type BaseResponse struct {
	Status       string `json:"status"`
	ResponseCode string `json:"response_code"`
	Message      string `json:"message"`
	Errors       string `json:"errors,omitempty"`
}

type GeneralResponse[T any] struct {
	BaseResponse
	Data T `json:"transaction_data,omitempty"`
}
