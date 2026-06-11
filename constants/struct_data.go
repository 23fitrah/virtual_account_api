package constants

type DataMaintenance struct {
	UserApprove string
	Action      string
	TglApprove  string
}

type DataSts struct {
	CloseDate    string
	KdStatus     string
	CloseRemarks string
	ProcessDate  string
	NamaAsliTrx  string
	BenefName    string
	Remark       string

	PassCheckCharges string
	Charges          string
	ChargesAmendment string
	ReffTracerAmend  string
}

type DataGPI struct {
	RowID                     string
	Reff                      string
	KdStatus                  string
	Status                    string
	Reason                    string
	Amount                    string
	TransactionIdentification string
	Currency                  string
	TransactionDate           string
	CrAmount                  string
	CrCurrency                string
	Rate                      string
	Charges                   string
	Mxtype                    string
}
