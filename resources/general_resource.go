package resources

type BaseResponse struct {
	ResponseId   string `json:"responseId"`
	ErrorCode    string `json:"errorCode"`
	StatusCode   string `json:"statusCode"`
	StatusDesc   string `json:"statusDesc"`
	ResponseTime string `json:"responseTime"`
}

type GeneralResponse[T any] struct {
	BaseResponse
	Data T `json:"transactionData,omitempty"`
}
