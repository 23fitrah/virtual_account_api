package constants

const (
	// Succes Code VA
	CodeVaCreate    = "VA_CREATED"
	CodeVaFound     = "VA_FOUND"
	CodeVaGetStatus = "VA_GET_STATUS"
	CodeVaGet       = "VA_GET_SUCCESS"
	//Success Code Payment
	CodePaymentReceived         = "PAYMENT_SUCCESS"
	CodePaymentAlreadyProcessed = "PAYMENT_ALREADY_PROCESSED"

	// Error Code VA
	CodeVaNotFound      = "VA_NOT_FOUND"
	CodeVaAlreadyExists = "VA_ALREADY_EXISTS"
	CodeVaInactive      = "VA_INACTIVE"
	CodeVaExpired       = "VA_EXPIRED"
	CodeVaCancelled     = "VA_CANCELLED"
	CodeVaAlreadyPaid   = "VA_ALREADY_PAID"

	// Error Merchant
	CodeMerchantNotFound  = "MERCHANT_NOT_FOUND"
	CodeMerchantInactive  = "MERCHANT_INACTIVE"
	CodeDuplicateMerchant = "DUPLICATE_CUSTOMER_REFERENCE"

	//Error Payment
	CodePaymentFailed    = "PAYMENT_FAILED"
	CodeAmountMismatch   = "AMOUNT_MISMATCH"
	CodePaymentInvalid   = "INVALID_PAYMENT_STATUS"
	CodeDuplicatePayment = "PAYMENT_REFERENCE_DUPLICATE"

	//Error General
	CodeAuthenticationFailed = "AUTHENTICATION_FAILED"
	CodeAuthorizationFailed  = "AUTHORIZATION_FAILED"
	CodeBadRequest           = "BAD_REQUEST"
	CodeInternalServerError  = "INTERNAL_SERVER_ERROR"
	CodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	CodeTransactionFailed    = "TRANSACTION_FAILED"

	CodeFailedToProcess  = "0004"
	CodeErrorSendMidTier = "0002"

	CodeTransactionSuccess = "00"
)
