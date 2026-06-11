package tests

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"virtual_account_api/internal/handlers"
	"virtual_account_api/internal/injector"
	"virtual_account_api/internal/repositories"
	"virtual_account_api/internal/routes"
	"virtual_account_api/internal/services"

	"github.com/alicebob/miniredis/v2"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// --- Test Models (Replica of Database Tables) ---

type AdapterServiceUser struct {
	Username string `gorm:"column:USERNAME"`
	Password string `gorm:"column:PASSWORD"`
}

func (AdapterServiceUser) TableName() string { return "ADAPTERSERVICEUSER" }

type SwiftGetKurs struct {
	Currency       string `gorm:"column:CURRENCY"`
	DebetBookRate  string `gorm:"column:DEBETBOOKRATE"`
	CreditBookRate string `gorm:"column:CREDITBOOKRATE"`
}

func (SwiftGetKurs) TableName() string { return "SWIFTGETKURS" }

type BriStatusCode struct {
	StatusCode  string `gorm:"column:StatusCode"`
	Description string `gorm:"column:Description"`
}

func (BriStatusCode) TableName() string { return "BRISTATUSCODE" }

type Bri3rdPty struct {
	BIC          string `gorm:"column:BIC"`
	CurrTrx      string `gorm:"column:curr_trx"`
	MTType       string `gorm:"column:mttype"`
	CancelMTType string `gorm:"column:cancel_mt_type"`
	CancelFee    string `gorm:"column:cancel_fee"`
	NosVos       string `gorm:"column:nos_vos"`
}

func (Bri3rdPty) TableName() string { return "BRI3RDPTY" }

type BriSwiftMaintenanceMX struct {
	ID          string    `gorm:"column:ID"`
	Status      string    `gorm:"column:STATUS"`
	Reff        string    `gorm:"column:REFF"`
	RowID       string    `gorm:"column:ROWID"`
	UserApprove string    `gorm:"column:USERAPPROVE"`
	Action      string    `gorm:"column:ACTION"`
	TglApprove  time.Time `gorm:"column:TGLAPPROVE"`
}

func (BriSwiftMaintenanceMX) TableName() string { return "BRISWIFTMAINTENANCEMX" }

type BriSwiftStsTrxMx struct {
	Reff             string    `gorm:"column:REFF"`
	RowID            string    `gorm:"column:ROW_ID"`
	KdStatus         string    `gorm:"column:KD_STATUS"`
	ErrDesc          string    `gorm:"column:ERR_DESC"`
	CloseDate        time.Time `gorm:"column:CLOSEDATE"`
	CloseRemarks     string    `gorm:"column:CLOSEREMARKS"`
	ProcessDate      time.Time `gorm:"column:PROCESS_DATE"`
	CountryOutgoing  string    `gorm:"column:COUNTRY_OUTGOING"`
	BicOutgoing      string    `gorm:"column:BIC_OUTGOING"`
	NamaBrinetsFull  string    `gorm:"column:NAMA_BRINETS_FULL"`
	NamaBrinets      string    `gorm:"column:NAMA_BRINETS"`
	AmtTrx           string    `gorm:"column:AMT_TRX"`
	Curr             string    `gorm:"column:CURR"`
	BenefName        string    `gorm:"column:BENEFNAME"`
	RejectCode       string    `gorm:"column:REJECT_CODE"`
	RejectDesc       string    `gorm:"column:REJECT_DESC"`
	SHA              string    `gorm:"column:SHA"`
	Remark2          string    `gorm:"column:REMARK2"`
	Charges          string    `gorm:"column:CHARGES"`
	ChargesAmendment string    `gorm:"column:CHARGESAMENDMENT"`
	PassCheckCharges string    `gorm:"column:PASSCHECKCHARGES"`
	ReffTracerAmend  string    `gorm:"column:REFF_TRACER_AMEND"`
	TraceCounter     string    `gorm:"column:TRACE_COUNTER"`
}

func (BriSwiftStsTrxMx) TableName() string { return "BRISWIFTSTSTRXMX" }

type BriSwiftTransDataTrxMx struct {
	SenderBIC            string    `gorm:"column:SENDERBIC"`
	RowID                string    `gorm:"column:ROWID"`
	Amount               string    `gorm:"column:AMOUNT"`
	Currency             string    `gorm:"column:CURRENCY"`
	ValueDate            time.Time `gorm:"column:VALUEDATE"`
	BenefAccount         string    `gorm:"column:BENEFACCOUNT"`
	BenefName            string    `gorm:"column:BENEFNAME"`
	MXType               string    `gorm:"column:MXTYPE"`
	MX                   string    `gorm:"column:MX"`
	CoverPlace           string    `gorm:"column:COVERPLACE"`
	ProcessDate          string    `gorm:"column:PROCESS_DATE"`
	InstrIdPain001       string    `gorm:"column:INSTRID_PAIN001"`
	InstrID              string    `gorm:"column:INSTRID"`
	RmtInf               string    `gorm:"column:RMTINF"`
	UETR                 string    `gorm:"column:UETR"`
	Amount32A            string    `gorm:"column:AMOUNT32A"`
	InstgRmbrsmntAgtAcct string    `gorm:"column:INSTGRMBRSMNTAGTACCT"`
	InstdRmbrsmntAgt     string    `gorm:"column:INSTDRMBRSMNTAGT"`
	IntrMyAgt1Acct       string    `gorm:"column:INTRMYAGT1ACCT"`
	SHA                  string    `gorm:"column:SHA"`
}

func (BriSwiftTransDataTrxMx) TableName() string { return "BRISWIFTTRANSDATATRXMX" }

type BriSwiftGpi struct {
	RowID                     string `gorm:"column:ROWID;primaryKey"`
	Reff                      string `gorm:"column:REFF"`
	KdStatus                  string `gorm:"column:KD_STATUS"`
	From                      string `gorm:"column:FROM_"`
	BusinessService           string `gorm:"column:BUSINESS_SERVICE"`
	TransactionIdentification string `gorm:"column:TRANSACTION_IDENTIFICATION"`
	InstructionIdentification string `gorm:"column:INSTRUCTION_IDENTIFICATION"`
	Originator                string `gorm:"column:ORIGINATOR"`
	Status                    string `gorm:"column:STATUS"`
	Reason                    string `gorm:"column:REASON"`
	Amount                    string `gorm:"column:AMOUNT"`
	Currency                  string `gorm:"column:CURRENCY"`
	TransactionDate           string `gorm:"column:TRANSACTION_DATE"`
	CrAmount                  string `gorm:"column:CRAMOUNT"`
	CrCurr                    string `gorm:"column:CRCURR"`
	Rate                      string `gorm:"column:RATE"`
	Charges                   string `gorm:"column:CHARGES"`
	SendSts                   string `gorm:"column:SEND_STS"`
	SendDesc                  string `gorm:"column:SEND_DESC"`
	FundAvailable             string `gorm:"column:FUND_AVAILABLE"`
	Kekinian                  string `gorm:"column:KEKINIAN"`
}

func (BriSwiftGpi) TableName() string { return "BRISWIFTGPI" }

type BriSwiftPembukuan struct {
	ID            int       `gorm:"column:ID;primaryKey"`
	RowID         string    `gorm:"column:ROWID"`
	TglLastProses time.Time `gorm:"column:TGLLASTPROSES"`
	Remark2       string    `gorm:"column:REMARK2"`
}

func (BriSwiftPembukuan) TableName() string { return "BRISWIFTPEMBUKUAN_K" }

type BriNostro struct {
	Bic     string `gorm:"column:BIC"`
	CurrTrx string `gorm:"column:CURR_TRX"`
	Nostro  string `gorm:"column:NOSTRO"`
}

func (BriNostro) TableName() string { return "BRINOSTRO" }

type ChargeTierNK struct {
	ID       string `gorm:"column:ID"`
	CurrTrx  string `gorm:"column:CURRTRX"`
	MinValue string `gorm:"column:MINVALUE"`
	Charge   string `gorm:"column:CHARGE"`
}

func (ChargeTierNK) TableName() string { return "CHARGETIERNK" }

// --- Setup ---

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB, *injector.AppContainer) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to sqlite: %v", err)
	}

	err = db.AutoMigrate(
		&AdapterServiceUser{},
		&SwiftGetKurs{},
		&BriStatusCode{},
		&Bri3rdPty{},
		&BriSwiftMaintenanceMX{},
		&BriSwiftStsTrxMx{},
		&BriSwiftTransDataTrxMx{},
		&BriSwiftGpi{},
		&BriSwiftPembukuan{},
		&BriNostro{},
		&ChargeTierNK{},
	)
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	// Setup mock redis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	t.Cleanup(func() {
		mr.Close()
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	virtualAccountRepo := repositories.NewVirtualAccountRepository()
	//cfg := config.configVa
	virtualAccountService := services.NewVirtualAccountService(db, redisClient, virtualAccountRepo, cfg)

	virtualAccountHandler := handlers.NewVirtualAccountHandler(virtualAccountService)

	container := &injector.AppContainer{
		DB:                    db,
		VirtualAccountHandler: virtualAccountHandler,
	}

	r := gin.New()
	api := r.Group("/api")
	routes.RegisterV1Router(api, container)

	return r, db, container
}

func TestV1Routes(t *testing.T) {
	password := "password"

	hash := md5.Sum([]byte(password))
	md5Hex := strings.ToUpper(hex.EncodeToString(hash[:]))

	mockUser := AdapterServiceUser{
		Username: "testuser",
		Password: md5Hex,
	}

	addAuth := func(body map[string]interface{}) map[string]interface{} {
		body["username"] = mockUser.Username
		body["password"] = password
		return body
	}

	t.Run("POST /api/v1/detail-transaction - Success", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)

		db.Create(&mockUser)
		db.Create(&BriStatusCode{StatusCode: "001", Description: "Success"})
		db.Create(&Bri3rdPty{BIC: "SENDERBIC", CurrTrx: "IDR", MTType: "103", CancelMTType: "C", CancelFee: "0", NosVos: "N"})

		// Use MX tables for IRK/SAK
		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW1", Amount: "1000", Currency: "IDR", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF123", InstrIdPain001: "PAIN1",
		})
		db.Create(&BriSwiftStsTrxMx{RowID: "ROW1", Reff: "REF123", KdStatus: "001"})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF123",
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("POST /api/v1/detail-transaction - Invalid Body", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		jsonBody := []byte(`{"invalid_json": `)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /api/v1/detail-transaction - Invalid Source", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"refNo":             "REF123",
			"transactionSource": "INVALID",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /api/v1/detail-transaction - Not Found", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"refNo":             "UNKNOWN",
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("POST /api/v1/detail-transaction - Authentication Failure", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := map[string]interface{}{
			"username":          mockUser.Username,
			"password":          "wrongpassword",
			"refNo":             "REF123",
			"transactionSource": "IRK",
		}
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRK) - TrxType 0006", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&BriStatusCode{StatusCode: "0001", Description: "Success"})
		db.Create(&Bri3rdPty{BIC: "SENDERBIC", CurrTrx: "USD", MTType: "103", CancelMTType: "C", CancelFee: "0", NosVos: "N"})

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRK_0006", Amount: "1000", Currency: "USD", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRK_0006",
		})
		db.Create(&BriSwiftStsTrxMx{
			RowID: "ROW_IRK_0006", Reff: "REF_IRK_0006", KdStatus: "0001",
			CountryOutgoing: "US", BicOutgoing: "OTHERBIC",
		})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRK_0006",
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"TrxType":"0006"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRK) - TrxType 0006 (ID + USD)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&BriStatusCode{StatusCode: "0001", Description: "Success"})
		db.Create(&Bri3rdPty{BIC: "SENDERBIC", CurrTrx: "USD", MTType: "103", CancelMTType: "C", CancelFee: "0", NosVos: "N"})

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRK_0006_ID", Amount: "1000", Currency: "USD", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRK_0006_ID",
		})
		// CountryOutgoing is ID, but Currency is USD. Should be 0006.
		db.Create(&BriSwiftStsTrxMx{
			RowID: "ROW_IRK_0006_ID", Reff: "REF_IRK_0006_ID", KdStatus: "0001",
			CountryOutgoing: "ID", BicOutgoing: "OTHERBIC",
		})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRK_0006_ID",
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"TrxType":"0006"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRK) - TrxType 0007", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&BriStatusCode{StatusCode: "0001", Description: "Success"})
		db.Create(&Bri3rdPty{BIC: "SENDERBIC", CurrTrx: "IDR", MTType: "103", CancelMTType: "C", CancelFee: "0", NosVos: "N"})

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRK_2", Amount: "1000", Currency: "IDR", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRK_2",
		})
		db.Create(&BriSwiftStsTrxMx{
			RowID: "ROW_IRK_2", Reff: "REF_IRK_2", KdStatus: "0001",
			CountryOutgoing: "ID", BicOutgoing: "OTHERBIC",
		})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRK_2",
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"TrxType":"0007"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRK) - TrxType 0007 with Status 44444", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&BriStatusCode{StatusCode: "44444", Description: "Success"})
		db.Create(&Bri3rdPty{BIC: "SENDERBIC", CurrTrx: "IDR", MTType: "103", CancelMTType: "C", CancelFee: "0", NosVos: "N"})

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRK_44444", Amount: "1000", Currency: "IDR", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRK_44444",
		})
		db.Create(&BriSwiftStsTrxMx{
			RowID: "ROW_IRK_44444", Reff: "REF_IRK_44444", KdStatus: "44444",
			CountryOutgoing: "ID", BicOutgoing: "OTHERBIC",
		})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRK_44444",
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"TrxType":"0007"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRK) - TrxType 0005 Fallback", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&BriStatusCode{StatusCode: "0002", Description: "Success"})
		db.Create(&Bri3rdPty{BIC: "SENDERBIC", CurrTrx: "IDR", MTType: "103", CancelMTType: "C", CancelFee: "0", NosVos: "N"})

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRK_0005", Amount: "1000", Currency: "IDR", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRK_0005",
		})
		// Condition for 0007: CountryOutgoing == "ID" && Currency == "IDR" && BicOutgoing != "BRINIDJA"
		// But KdStatus is NOT "0001" or "44444" -> Should fall back to 0005
		db.Create(&BriSwiftStsTrxMx{
			RowID: "ROW_IRK_0005", Reff: "REF_IRK_0005", KdStatus: "0002",
			CountryOutgoing: "ID", BicOutgoing: "OTHERBIC",
		})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRK_0005",
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"TrxType":"0005"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRK) - TrxType 0005 Fallback (BicOutgoing is BRINIDJA)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&BriStatusCode{StatusCode: "0001", Description: "Success"})
		db.Create(&Bri3rdPty{BIC: "SENDERBIC", CurrTrx: "IDR", MTType: "103", CancelMTType: "C", CancelFee: "0", NosVos: "N"})

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRK_BRINIDJA", Amount: "1000", Currency: "IDR", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRK_BRINIDJA",
		})
		// Condition for 0007 is NOT met because BicOutgoing is "BRINIDJA"
		db.Create(&BriSwiftStsTrxMx{
			RowID:           "ROW_IRK_BRINIDJA",
			Reff:            "REF_IRK_BRINIDJA",
			KdStatus:        "0001",
			CountryOutgoing: "ID",
			BicOutgoing:     "BRINIDJA",
		})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRK_BRINIDJA",
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"TrxType":"0005"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRN)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		// Use MX tables for IRN/SANK
		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRN", Amount: "2000", Currency: "USD", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRN",
		})
		db.Create(&BriSwiftStsTrxMx{RowID: "ROW_IRN", Reff: "REF_IRN", KdStatus: "001", CountryOutgoing: "US"})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRN",
			"transactionSource": "IRN",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRN) - SenderBicIsDepcor & TrxType 0003", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRN_2", Amount: "2000", Currency: "IDR", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRN_2", CoverPlace: "SENDERBIC",
		})
		db.Create(&BriSwiftStsTrxMx{RowID: "ROW_IRN_2", Reff: "REF_IRN_2", KdStatus: "001", CountryOutgoing: "ID"})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRN_2",
			"transactionSource": "IRN",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"SenderBicIsDepcor":"1"`)
		assert.Contains(t, w.Body.String(), `"TrxType":"0003"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRN) - TrxType 0002", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRN_0002", Amount: "2000", Currency: "USD", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRN_0002",
		})
		db.Create(&BriSwiftStsTrxMx{RowID: "ROW_IRN_0002", Reff: "REF_IRN_0002", KdStatus: "001", CountryOutgoing: "US"})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRN_0002",
			"transactionSource": "IRN",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"TrxType":"0002"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRN) - Default MXType", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRN_MX", Amount: "2000", Currency: "USD", ValueDate: time.Now(),
			MXType: "", InstrID: "REF_IRN_MX", // Empty MTType
		})
		db.Create(&BriSwiftStsTrxMx{RowID: "ROW_IRN_MX", Reff: "REF_IRN_MX", KdStatus: "001", CountryOutgoing: "US"})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRN_MX",
			"transactionSource": "IRN",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"MXtype":"103"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRN) - With Nostro", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&BriNostro{Bic: "COVER123XXX", CurrTrx: "USD", Nostro: "1234567890"})

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRN_NOSTRO", Amount: "2000", Currency: "USD", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRN_NOSTRO", CoverPlace: "COVER123",
		})
		db.Create(&BriSwiftStsTrxMx{RowID: "ROW_IRN_NOSTRO", Reff: "REF_IRN_NOSTRO", KdStatus: "001", CountryOutgoing: "US"})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRN_NOSTRO",
			"transactionSource": "IRN",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"NostroAccount":"1234567890"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRN) - SenderBicIsDepcor False", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRN_ND", Amount: "2000", Currency: "USD", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRN_ND", CoverPlace: "OTHERBIC",
		})
		db.Create(&BriSwiftStsTrxMx{RowID: "ROW_IRN_ND", Reff: "REF_IRN_ND", KdStatus: "001", CountryOutgoing: "US"})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRN_ND",
			"transactionSource": "IRN",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"SenderBicIsDepcor":"0"`)
	})

	t.Run("POST /api/v1/detail-transaction - Success (IRN) - TrxType 0001", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRN_0001", Amount: "2000", Currency: "IDR", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_IRN_0001",
		})
		// CountryOutgoing is "US" (not empty), but Currency is "IDR" -> TrxType 0001 (because 0002 requires Currency != IDR)
		db.Create(&BriSwiftStsTrxMx{RowID: "ROW_IRN_0001", Reff: "REF_IRN_0001", KdStatus: "001", CountryOutgoing: "US"})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_IRN_0001",
			"transactionSource": "IRN",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"TrxType":"0001"`)
	})

	t.Run("POST /api/v1/detail-transaction - IRK No 3rd Party", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		// No BRI3RDPTY created

		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_NO_3RD", Amount: "1000", Currency: "IDR", ValueDate: time.Now(),
			MXType: "103", InstrID: "REF_NO_3RD",
		})
		db.Create(&BriSwiftStsTrxMx{RowID: "ROW_NO_3RD", Reff: "REF_NO_3RD", KdStatus: "001"})

		body := addAuth(map[string]interface{}{
			"refNo":             "REF_NO_3RD",
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/detail-transaction", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Expect 404 because of INNER JOIN
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("POST /api/v1/incoming - Invalid Date", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"bic":             "TESTBIC",
			"transactionDate": "invalid-date",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/incoming", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /api/v1/incoming - Success (NCHG)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		rowID := "ROW_NCHG"
		db.Create(&BriSwiftTransDataTrxMx{SenderBIC: "TESTBIC", RowID: rowID, InstrID: "TAG" + rowID})
		db.Create(&BriSwiftStsTrxMx{RowID: rowID})
		db.Create(&BriSwiftPembukuan{
			ID: 100, RowID: rowID,
			TglLastProses: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			Remark2:       "Some CHG_ Remark",
		})

		body := addAuth(map[string]interface{}{
			"bic":             "TESTBIC",
			"transactionDate": "2023-01-01",
		})
		jsonBody, _ := sonic.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/incoming", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"transaction_code":"NCHG"`)
	})

	t.Run("POST /api/v1/incoming - Success (NTRF)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		rowID := "ROW_NTRF"
		db.Create(&BriSwiftTransDataTrxMx{SenderBIC: "TESTBIC", RowID: rowID, InstrID: "TAG" + rowID})
		db.Create(&BriSwiftStsTrxMx{RowID: rowID})
		db.Create(&BriSwiftPembukuan{
			ID: 101, RowID: rowID,
			TglLastProses: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			Remark2:       "Some Normal Remark",
		})

		body := addAuth(map[string]interface{}{
			"bic":             "TESTBIC",
			"transactionDate": "2023-01-01",
		})
		jsonBody, _ := sonic.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/incoming", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"transaction_code":"NTRF"`)
	})

	t.Run("POST /api/v1/incoming - Success (No Data)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"bic":             "TESTBIC",
			"transactionDate": "2023-01-01",
		})
		jsonBody, _ := sonic.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/incoming", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"totalData":"0"`)
	})

	t.Run("POST /api/v1/kurs-adapter - Success (Same Currency)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX001",
			"debetCurrency":  "USD",
			"debetAmount":    "100",
			"creditCurrency": "USD",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Get Currency Success")
	})

	t.Run("POST /api/v1/kurs-adapter - Success (USD to IDR)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		// Assuming a scale factor of 10,000,000. 14000 * 10,000,000 = 140000000000
		db.Create(&SwiftGetKurs{Currency: "IDR", DebetBookRate: "140000000000", CreditBookRate: "140000000000"})

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX002",
			"debetCurrency":  "USD",
			"debetAmount":    "100",
			"creditCurrency": "IDR",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Get Currency Success")
	})

	t.Run("POST /api/v1/kurs-adapter - Success (IDR to USD)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		// Assuming a scale factor of 10,000,000. 14000 * 10,000,000 = 140000000000
		db.Create(&SwiftGetKurs{Currency: "IDR", DebetBookRate: "140000000000", CreditBookRate: "140000000000"})

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX003",
			"debetCurrency":  "IDR",
			"debetAmount":    "1400000",
			"creditCurrency": "USD",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Get Currency Success")
	})

	t.Run("POST /api/v1/kurs-adapter - Not Found", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX004",
			"debetCurrency":  "USD",
			"debetAmount":    "100",
			"creditCurrency": "EUR",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Cannot find currency data")
	})

	t.Run("POST /api/v1/kurs-adapter - Invalid Amount", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX005",
			"debetCurrency":  "USD",
			"debetAmount":    "-100",
			"creditCurrency": "IDR",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /api/v1/kurs-adapter - Zero Amount", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX_ZERO",
			"debetCurrency":  "USD",
			"debetAmount":    "0",
			"creditCurrency": "IDR",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /api/v1/kurs-adapter - Success (EUR to IDR)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		// 1 USD = 0.9 EUR -> Rate = 0.9 * 10,000,000 = 9000000
		db.Create(&SwiftGetKurs{Currency: "EUR", DebetBookRate: "9000000", CreditBookRate: "9000000"})
		// 1 USD = 14000 IDR -> Rate = 14000 * 10,000,000 = 140000000000
		db.Create(&SwiftGetKurs{Currency: "IDR", DebetBookRate: "140000000000", CreditBookRate: "140000000000"})

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX006",
			"debetCurrency":  "EUR",
			"debetAmount":    "90",
			"creditCurrency": "IDR",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// 90 EUR / 0.9 = 100 USD. 100 USD * 14000 = 1,400,000 IDR
		assert.Contains(t, w.Body.String(), "1400000.00")
	})

	t.Run("POST /api/v1/kurs-adapter - Invalid Rate", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&SwiftGetKurs{Currency: "JPY", DebetBookRate: "0", CreditBookRate: "0"})

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX007",
			"debetCurrency":  "JPY",
			"debetAmount":    "1000",
			"creditCurrency": "USD",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("POST /api/v1/kurs-adapter - Debet Currency Not Found", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX008",
			"debetCurrency":  "UNKNOWN",
			"debetAmount":    "100",
			"creditCurrency": "USD",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("POST /api/v1/kurs-adapter - Credit Currency Invalid Rate", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&SwiftGetKurs{Currency: "JPY", DebetBookRate: "100", CreditBookRate: "0"})

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX009",
			"debetCurrency":  "USD",
			"debetAmount":    "100",
			"creditCurrency": "JPY",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("POST /api/v1/kurs-adapter - Invalid Amount Format", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX010",
			"debetCurrency":  "USD",
			"debetAmount":    "invalid",
			"creditCurrency": "IDR",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /api/v1/kurs-adapter - Small Amount", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		// 1 USD = 14000 IDR.
		db.Create(&SwiftGetKurs{Currency: "IDR", DebetBookRate: "140000000000", CreditBookRate: "140000000000"})

		// 1 IDR to USD = 0.0000714... -> 0.00
		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX_SMALL",
			"debetCurrency":  "IDR",
			"debetAmount":    "1",
			"creditCurrency": "USD",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"creditAmount":"0.00"`)
	})

	t.Run("POST /api/v1/kurs-adapter - Success (Comma in Amount)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionID":  "TRX_COMMA",
			"debetCurrency":  "USD",
			"debetAmount":    "100,50",
			"creditCurrency": "USD",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/kurs-adapter", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"debetAmount":"100,50"`)
		assert.Contains(t, w.Body.String(), `"creditAmount":"100.50"`)
	})

	t.Run("POST /api/v1/mt199/get-pending - Success", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&BriStatusCode{StatusCode: "KD_STATUS1", Description: "DESC1"})
		db.Create(&Bri3rdPty{BIC: "SENDERBIC", CurrTrx: "IDR", MTType: "103"})
		db.Create(&BriSwiftMaintenanceMX{Reff: "TAG20_1", Status: "4"})
		db.Create(&BriSwiftStsTrxMx{Reff: "TAG20_1", KdStatus: "KD_STATUS1", CountryOutgoing: "ID", BicOutgoing: "BIC1"})
		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW1", Amount: "1000", Currency: "IDR", ValueDate: time.Now(),
			MXType: "103", InstrID: "TAG20_1",
		})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/get-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("POST /api/v1/mt199/get-pending - Invalid Source", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionSource": "INVALID",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/get-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /api/v1/mt199/get-pending - Success (IRN)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&BriSwiftMaintenanceMX{Reff: "TAG20_IRN", Status: "4"})
		db.Create(&BriSwiftStsTrxMx{Reff: "TAG20_IRN", KdStatus: "KD_STATUS1"})
		db.Create(&BriSwiftTransDataTrxMx{
			SenderBIC: "SENDERBIC", RowID: "ROW_IRN", Amount: "1000", Currency: "USD", ValueDate: time.Now(),
			MXType: "103", InstrID: "TAG20_IRN",
		})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRN",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/get-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("POST /api/v1/mt199/get-pending - Get All with 100 data", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)
		db.Create(&Bri3rdPty{
			BIC:     "PAGINATION_BIC",
			CurrTrx: "IDR",
			MTType:  "103",
		})

		// Create 100 records
		for i := 0; i < 100; i++ {
			ref := "TAG20_P_" + strconv.Itoa(i)
			rowID := "ROW_P_" + strconv.Itoa(i)
			db.Create(&BriSwiftMaintenanceMX{ID: "P_" + strconv.Itoa(i), Reff: ref, RowID: rowID, Status: "4"})
			db.Create(&BriSwiftStsTrxMx{RowID: rowID, Reff: ref, KdStatus: "11000", CountryOutgoing: "ID", BicOutgoing: "BIC1"})
			db.Create(&BriSwiftTransDataTrxMx{
				SenderBIC: "PAGINATION_BIC", RowID: "ROW_P_" + strconv.Itoa(i), Amount: "1000", Currency: "IDR", ValueDate: time.Now(),
				MXType: "103", InstrID: ref,
			})
		}

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/get-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"totalData":"100"`)
	})

	t.Run("POST /api/v1/mt199/get-pending - No Data", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/get-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"totalData":"0"`)
	})

	t.Run("POST /api/v1/mt199/close-pending - Success", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		// Use MX tables
		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_ID", InstrID: "REFF123"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_ID", Reff: "REFF123", KdStatus: "11999", AmtTrx: "100", Curr: "IDR", BenefName: "BENEF"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_ID", Reff: "REFF123", RowID: "TRX_ID", Status: "9"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "TRX_ID",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/close-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
	})

	t.Run("POST /api/v1/mt199/close-pending - Data Not Found", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		// Use MX tables
		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_ID", InstrID: "REFF123"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_ID", Reff: "REFF123", KdStatus: "999", AmtTrx: "100", Curr: "IDR", BenefName: "BENEF"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_ID", Reff: "REFF123", RowID: "TRX_ID", Status: "5"})
		db.Create(&BriStatusCode{StatusCode: "999", Description: "Some Status"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "TRX_ID",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/close-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Not allowed")
	})

	t.Run("POST /api/v1/mt199/close-pending - Invalid Source", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionSource": "INVALID",
			"transactionID":     "TRX_ID",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/close-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /api/v1/mt199/close-pending - Real Not Found", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "UNKNOWN_ID_CLOSE",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/close-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("POST /api/v1/mt199/close-pending - Status Description Check", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriStatusCode{StatusCode: "11999", Description: "Request Approval - Rejected by System"})

		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_STAT", InstrID: "REF_STAT"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_STAT", Reff: "REF_STAT", KdStatus: "11999"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_STAT", Reff: "REF_STAT", RowID: "TRX_STAT", Status: "4"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "TRX_STAT",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/close-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Rejected by System")
	})

	t.Run("POST /api/v1/mt199/close-pending - Success (IRN)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_IRN", InstrID: "REF_IRN"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_IRN", Reff: "REF_IRN", KdStatus: "11999"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_IRN", Reff: "REF_IRN", RowID: "TRX_IRN", Status: "9"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRN",
			"transactionID":     "TRX_IRN",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/close-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
	})

	t.Run("POST /api/v1/mt199/close-pending - SANK GPI Exclusion (Status 92229)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		// SANK transaction, MXType 103, but KdStatus 92229
		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_SANK_GPI", InstrID: "REF_SANK_GPI", MXType: "103"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_SANK_GPI", Reff: "REF_SANK_GPI", KdStatus: "92229"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_SANK_GPI", Reff: "REF_SANK_GPI", RowID: "TRX_SANK_GPI", Status: "9"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRN",
			"transactionID":     "TRX_SANK_GPI",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/close-pending", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")

		// Verify GPI table is empty
		var count int64
		db.Model(&BriSwiftGpi{}).Count(&count)
		assert.Equal(t, int64(0), count, "GPI table should be empty for SANK with status 92229")
	})

	t.Run("POST /api/v1/mt199/reject - Success", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_REJECT", InstrID: "REF_REJECT"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_REJECT", Reff: "REF_REJECT", KdStatus: "11000"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_REJECT", Reff: "REF_REJECT", RowID: "TRX_REJECT", Status: "4", UserApprove: "user1"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "TRX_REJECT",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/reject", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
	})

	t.Run("POST /api/v1/mt199/reject - Not Found", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "UNKNOWN_ID",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/reject", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("POST /api/v1/mt199/reject - Already Processed", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		// Status 9 usually means closed/approved, so GetDataMt199 (pending) should return nil
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_PROC", Reff: "REF_PROC", RowID: "TRX_PROC", Status: "9"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "TRX_PROC",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/reject", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("POST /api/v1/mt199/reject - Success (IRN)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_REJECT_IRN", InstrID: "REF_REJECT_IRN"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_REJECT_IRN", Reff: "REF_REJECT_IRN", KdStatus: "11000"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_REJECT_IRN", Reff: "REF_REJECT_IRN", RowID: "TRX_REJECT_IRN", Status: "4", UserApprove: "user1"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRN",
			"transactionID":     "TRX_REJECT_IRN",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/reject", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
	})

	t.Run("POST /api/v1/mt199/flag - Success", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_FLAG", InstrID: "REF_FLAG"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_FLAG", Reff: "REF_FLAG", KdStatus: "11000"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_FLAG", Reff: "REF_FLAG", RowID: "TRX_FLAG", Status: "4", UserApprove: "user1"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "TRX_FLAG",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/flag", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
	})

	t.Run("POST /api/v1/mt199/flag - Not Found", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "UNKNOWN_ID",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/flag", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("POST /api/v1/mt199/flag - Already Processed", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_PROC_FLAG", Reff: "REF_PROC_FLAG", RowID: "TRX_PROC_FLAG", Status: "9"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "TRX_PROC_FLAG",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/flag", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("POST /api/v1/mt199/flag - Success (IRN)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_FLAG_IRN", InstrID: "REF_FLAG_IRN"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_FLAG_IRN", Reff: "REF_FLAG_IRN", KdStatus: "11000"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_FLAG_IRN", Reff: "REF_FLAG_IRN", RowID: "TRX_FLAG_IRN", Status: "4", UserApprove: "user1"})

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRN",
			"transactionID":     "TRX_FLAG_IRN",
			"userRequest":       "user1",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/flag", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
	})

	t.Run("POST /api/v1/mt199/release - Success", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_RELEASE", InstrID: "REF_RELEASE"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_RELEASE", Reff: "REF_RELEASE", KdStatus: "12999"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_RELEASE", Reff: "REF_RELEASE", RowID: "TRX_RELEASE", Status: "9", UserApprove: "user1"})

		body := addAuth(map[string]interface{}{
			"transactionSource":  "IRK",
			"transactionID":      "TRX_RELEASE",
			"userRequest":        "user1",
			"remarkNewBenefName": "New Benef",
			"traceCounter":       "123",
			"reffTracer":         "REF123",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/release", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
	})

	t.Run("POST /api/v1/mt199/release - Because data not found", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_RELEASE", InstrID: "REF_RELEASE"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_RELEASE", Reff: "REF_RELEASE", KdStatus: "0002"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_RELEASE", Reff: "REF_RELEASE", RowID: "TRX_RELEASE", Status: "4", UserApprove: "user1"})

		body := addAuth(map[string]interface{}{
			"transactionSource":  "IRK",
			"transactionID":      "TRX_RELEASE",
			"userRequest":        "user1",
			"remarkNewBenefName": "New Benef",
			"traceCounter":       "123",
			"reffTracer":         "REF123",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/release", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Not allowed to close transaction.")
	})

	t.Run("POST /api/v1/mt199/release - Not Found", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		body := addAuth(map[string]interface{}{
			"transactionSource": "IRK",
			"transactionID":     "UNKNOWN_ID",
			"userRequest":       "user1",
			"reffTracer":        "REF123",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/release", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /api/v1/mt199/release - Already Processed", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		// Exists in Status Trx but not in Maintenance (Waiting Confirmation)
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_REL_PROC", Reff: "REF_REL_PROC", KdStatus: "001"})

		body := addAuth(map[string]interface{}{
			"transactionSource":  "IRK",
			"transactionID":      "TRX_REL_PROC",
			"userRequest":        "user1",
			"remarkNewBenefName": "New Benef",
			"reffTracer":         "REF123",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/release", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /api/v1/mt199/release - Success (IRN)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_RELEASE_IRN", InstrID: "REF_RELEASE_IRN"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_RELEASE_IRN", Reff: "REF_RELEASE_IRN", KdStatus: "12999"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_RELEASE_IRN", Reff: "REF_RELEASE_IRN", RowID: "TRX_RELEASE_IRN", Status: "9", UserApprove: "user1"})

		body := addAuth(map[string]interface{}{
			"transactionSource":  "IRN",
			"transactionID":      "TRX_RELEASE_IRN",
			"userRequest":        "user1",
			"remarkNewBenefName": "New Benef",
			"traceCounter":       "123",
			"reffTracer":         "REF123",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/release", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
	})

	t.Run("POST /api/v1/mt199/release - Failed because data not found (IRN)", func(t *testing.T) {
		r, db, _ := setupTestRouter(t)
		db.Create(&mockUser)

		db.Create(&BriSwiftTransDataTrxMx{RowID: "TRX_RELEASE_IRN", InstrID: "REF_RELEASE_IRN"})
		db.Create(&BriSwiftStsTrxMx{RowID: "TRX_RELEASE_IRN", Reff: "REF_RELEASE_IRN", KdStatus: "0002"})
		db.Create(&BriSwiftMaintenanceMX{ID: "TRX_RELEASE_IRN", Reff: "REF_RELEASE_IRN", RowID: "TRX_RELEASE_IRN", Status: "4", UserApprove: "user1"})

		body := addAuth(map[string]interface{}{
			"transactionSource":  "IRN",
			"transactionID":      "TRX_RELEASE_IRN",
			"userRequest":        "user1",
			"remarkNewBenefName": "New Benef",
			"traceCounter":       "123",
			"reffTracer":         "REF123",
		})
		jsonBody, _ := sonic.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/mt199/release", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Not allowed to close transaction.")
	})
}
