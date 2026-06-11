package constants

const (
	// Context Keys for Audit Logging
	ContextKeyOldValue        = "oldValue"
	ContextKeyNewValue        = "newValue"
	ContextKeyResponseMessage = "responseMessage"
	ContextKeyUser            = "username"
	ContextKeyMenu            = "menu"

	// Success
	StatusGetSuccess         = "Get Data Success"
	StatusGetCurrencysuccess = "Get Currency Success"

	// Generic prefixes
	StatusPrefixSuccess = "[Success]"
	StatusPrefixFailed  = "[Failed]"
	StatusPrefixTimeout = "[Timeout]"

	// System / validation errors
	StatusIncompleteData                        = "[Failed] Incomplete Data"
	StatusInvalidDateFormat                     = "[Failed] Invalid Date Format"
	StatusInvalidCredential                     = "[Failed] Invalid Username or Password"
	StatusInvalidDBConnection                   = "[Timeout] Invalid DB Connection"
	StatusInvalidTransactionSource              = "[Failed] Invalid Transaction Source"
	StatusDataNotFound                          = "[Empty] Data Not Found"
	StatusErrorCustom                           = "[Failed] Failed : "
	StatusEndpointNotFound                      = "[Failed] Endpoint Not Found"
	StatusInvalidKursValue                      = "Invalid kurs value"
	StatusSuccessCloseData                      = "[Success] Transaction Successfully Closed"
	StatusSuccessRejectData                     = "[Success] Transaction Successfully Rejected"
	StatusSuccessFlagData                       = "[Success] Transaction Successfully Approved"
	StatusSuccessReleaseData                    = "[Success] Transaction Successfully Released"
	StatusSwiftNotAllowedCloseTransaction       = "Not allowed to close transaction. Current Swift Adapter Status : "
	StatusMaintenanceNotAllowedCloseTransaction = "Not allowed to close transaction. Current Swift Adapter Maintenance Status : "
	StatusFailedCloseData                       = "[Failed] Close data failed, please check status. "
	StatusFailedRejectData                      = "[Failed] Reject data failed, please check status. "
	StatusFailedFlagData                        = "[Failed] Approve data failed, please check status. "
	StatusFailedReleaseData                     = "[Failed] Approve data failed, please check status. "
	StatusInvalidRequestBody                    = "[Failed] Invalid request body"
	StatusInvalidAuthRequest                    = "[Failed] Invalid Username or Password"
	StatusInvalidDebetAmount                    = "[Exception] Failed : invalid debet amount"
	StatusInvalidConversionAmount               = "[Failed] Invalid conversion amount : "
	StatusCannotFindCurrencyData                = "[Failed] Cannot find currency data : "
	StatusInvalidDebetKurs                      = "[Exception] Failed : invalid debet kurs for "
	StatusFailedInsertHistoryLog                = "[FAILED] Failed to insert history log async"
	StatusFailedInsertOCServiceLog              = "[FAILED] Failed to insert oc-service log async"
	StatusFailedAudit                           = "[FAILED] Failed to insert audit"
	StatusWarningAuditLogFull                   = "[WARNING] Audit Log Buffer Full - Dropping HistoryLog"

	// Generic fallback
	StatusUnknown = "Unknown Status"
)
