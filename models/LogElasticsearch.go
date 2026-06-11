package models

type ElasticLog struct {
	Timestamp     string `json:"timestamp"`
	TransactionID string `json:"transaction_id"`
	FullRequest   string `json:"full_request"`
	FullResponse  string `json:"full_response"`
	Message       string `json:"message"`
	ResponseCode  int    `json:"response_code"`
	IpSource      string `json:"ip_source"`
	UserID        string `json:"user_id"`
	Function      string `json:"function"`
	URL           string `json:"url"`
	DurationMs    int64  `json:"duration_ms"`
}
