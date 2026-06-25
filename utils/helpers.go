package utils

import (
	"crypto/md5"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func GetModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Println("Warning: go.mod not found, using 'opphicf2-fase-2b'")
		return "opphicf2-fase-2b"
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	return "opphicf2-fase-2b"
}

func SlugifyFilename(original string) string {
	ext := ""
	dotIdx := strings.LastIndex(original, ".")
	if dotIdx != -1 {
		ext = original[dotIdx:]
		original = original[:dotIdx]
	}
	return slug.Make(original) + ext
}

func CamelToSnake(s string) string {
	re := regexp.MustCompile("(.)([A-Z][a-z]+)")
	s = re.ReplaceAllString(s, "${1}_${2}")
	re = regexp.MustCompile("([a-z0-9])([A-Z])")
	s = re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(s)
}

type JSONBMap map[string]interface{}

// Untuk membaca dari DB (Scan)
func (j *JSONBMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Scan source is not []byte")
	}
	return json.Unmarshal(bytes, j)
}

// Untuk menyimpan ke DB (Value)
func (j JSONBMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func WriteJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

func DateTimeNow() string {
	return time.Now().Format("020106150405")
}

func GetStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func ToTitleCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

func ToStringReplace(s string, delimiter string) string {
	return strings.ReplaceAll(s, delimiter, " ")
}

func notInArray(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return false
		}
	}
	return true
}

func ParseDecimal(val string) (decimal.Decimal, error) {
	val = strings.ReplaceAll(val, ",", ".")
	return decimal.NewFromString(val)
}

func FormatMaintenanceStatus(status string) string {
	if status == "" {
		return ""
	}

	switch status {
	case "0":
		return "0 - "
	case "1":
		return "1 - "
	case "2":
		return "2 - "
	case "4":
		return "4 - "
	case "41":
		return "4 - "
	case "5":
		return "5 - "
	case "6":
		return "6 - "
	case "9":
		return "9 - "
	default:
		return status + " - "
	}
}

func Md5Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func CalculateCharges(
	senderBIC string,
	sha string,
	traceCounter decimal.Decimal,
	baseCharge decimal.Decimal,
	baseChargeUSD decimal.Decimal,
) decimal.Decimal {
	if baseChargeUSD.IsZero() {
		return decimal.Zero
	}

	sha = strings.ToUpper(sha)

	localBankSwift := false
	if len(senderBIC) > 6 {
		if strings.ToUpper(senderBIC[4:6]) == "ID" {
			localBankSwift = true
		}
	}

	offset := decimal.Zero

	switch {
	case localBankSwift && sha == "OUR":
		offset = decimal.NewFromInt(4)
	case localBankSwift && sha != "OUR":
		offset = decimal.NewFromInt(5)
	case !localBankSwift && sha == "OUR":
		offset = decimal.NewFromInt(5)
	default:
		offset = decimal.NewFromInt(6)
	}

	return traceCounter.
		Add(offset).
		Mul(decimal.NewFromInt(5)).
		Mul(baseCharge).
		Div(baseChargeUSD)
}

func GetChargeSank(db *gorm.DB, curr string) int64 {
	var charge int64

	err := db.Table("CHARGETIERNK").Where("CURRTRX = ?", curr).Select("CHARGE").Scan(&charge).Error
	if err != nil {
		return 0
	}

	return charge
}

func ToString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func StructToJSONString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
