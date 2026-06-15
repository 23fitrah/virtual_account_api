package constants

const (
	// Succes Code VA
	StatusCodeVaCreate    = "VA_CREATED"
	StatusCodeVaFound     = "VA_FOUND"
	StatusCodeVaGetStatus = "VA_GET_STATUS"
	StatusCodeVaGet       = "VA_GET_SUCCESS"
	StatusCodeVaSuccess   = "VA_SUCCESS"
	StatusCodeVaFailed    = "VA_FAILED"
	CodeVaSuccess         = "VA-0000"
	CodeVaFailed          = "VA-0004"

	//Success Code Payment
	StatusCodePaymentReceived         = "PAYMENT_SUCCESS"
	StatusCodePaymentAlreadyProcessed = "PAYMENT_ALREADY_PROCESSED"

	// Error Code VA
	StatusCodeVaNotFound      = "VA_NOT_FOUND"
	StatusCodeVaAlreadyExists = "VA_ALREADY_EXISTS"
	StatusCodeVaInactive      = "VA_INACTIVE"
	StatusCodeVaExpired       = "VA_EXPIRED"
	StatusCodeVaCancelled     = "VA_CANCELLED"
	StatusCodeVaAlreadyPaid   = "VA_ALREADY_PAID"

	// Error Merchant
	StatusCodeMerchantNotFound  = "MERCHANT_NOT_FOUND"
	StatusCodeMerchantInactive  = "MERCHANT_INACTIVE"
	StatusCodeDuplicateMerchant = "DUPLICATE_CUSTOMER_REFERENCE"

	//Error Payment
	StatusCodePaymentFailed    = "PAYMENT_FAILED"
	StatusCodeAmountMismatch   = "AMOUNT_MISMATCH"
	StatusCodePaymentInvalid   = "INVALID_PAYMENT_STATUS"
	StatusCodeDuplicatePayment = "PAYMENT_REFERENCE_DUPLICATE"

	//Error General
	StatusCodeAuthenticationFailed = "AUTHENTICATION_FAILED"
	StatusCodeAuthorizationFailed  = "AUTHORIZATION_FAILED"
	StatusCodeBadRequest           = "BAD_REQUEST"
	StatusCodeInternalServerError  = "INTERNAL_SERVER_ERROR"
	StatusCodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	StatusCodeTransactionFailed    = "TRANSACTION_FAILED"

	StatuCodeFailedToProcess   = "VA-0004"
	StatusCodeErrorSendMidTier = "VA-0002"

	CodeTransactionSuccess = "VA-0000"

	CodeEndpointNotFound = "0004"
	CodeErrorSendMidTier = "0002"
)
