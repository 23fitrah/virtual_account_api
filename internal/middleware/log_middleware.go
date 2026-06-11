package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"
	"virtual_account_api/models"
	"virtual_account_api/utils"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type RequestExtractor struct {
	RefNo string `json:"refNo"`
}

func ESMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		trxID := ""

		var requestBodyBytes []byte
		if c.Request.Body != nil {
			requestBodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

		if len(requestBodyBytes) > 0 {
			var extractor RequestExtractor
			if err := sonic.Unmarshal(requestBodyBytes, &extractor); err == nil {
				if extractor.RefNo != "" {
					trxID = extractor.RefNo
				}
			}
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		duration := time.Since(startTime).Milliseconds()

		username := c.GetString("username")
		if username == "" {
			username = "guest"
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			loc = time.FixedZone("WIB", 7*60*60)
		}
		timestampWIB := startTime.In(loc).Format(time.RFC3339)

		logData := models.ElasticLog{
			Timestamp:     timestampWIB,
			TransactionID: trxID,
			FullRequest:   string(requestBodyBytes),
			FullResponse:  blw.body.String(),
			Message:       http.StatusText(c.Writer.Status()),
			ResponseCode:  c.Writer.Status(),
			IpSource:      c.ClientIP(),
			UserID:        username,
			Function:      c.HandlerName(),
			URL:           c.Request.URL.Path,
			DurationMs:    duration,
		}

		utils.LogJSON(logData)
	}
}
