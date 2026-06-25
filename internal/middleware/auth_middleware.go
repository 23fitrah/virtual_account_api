package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"strings"
	"virtual_account_api/constants"
	"virtual_account_api/internal/repositories"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthWSMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		const prefix = "Basic "
		if !strings.HasPrefix(authHeader, prefix) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": constants.CodeErrorSendMidTier,
				"status_desc": constants.StatusUnAuth,
			})
			c.Abort()
			return
		}

		// Decode base64 "username:password"
		encodedCreds := strings.TrimPrefix(authHeader, prefix)
		decoded, err := base64.StdEncoding.DecodeString(encodedCreds)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": constants.CodeErrorSendMidTier,
				"status_desc": constants.StatusInvalidAuthEncoding,
			})
			c.Abort()
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": constants.CodeErrorSendMidTier,
				"status_desc": constants.StatusInvalidAuthCredential,
			})
			c.Abort()
			return
		}

		username := credentials[0]
		password := credentials[1]

		// Ambil kredensial yang valid dari ENV (sama pola seperti TOKEN_AUTH)
		validUsername := os.Getenv("BASIC_AUTH_USERNAME")
		validPassword := os.Getenv("BASIC_AUTH_PASSWORD")

		if username != validUsername || password != validPassword {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": constants.CodeErrorSendMidTier,
				"status_desc": constants.StatusInvalidAuthCredential,
			})
			c.Abort()
			return
		}
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": constants.CodeErrorSendMidTier,
				"status_desc": constants.StatusInvalidRequestBody,
			})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := sonic.Unmarshal(bodyBytes, &payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": constants.CodeErrorSendMidTier,
				"status_desc": constants.StatusInvalidAuthRequest,
			})
			c.Abort()
			return
		}

		if payload.Username == "" || payload.Password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": constants.CodeErrorSendMidTier,
				"status_desc": constants.StatusInvalidAuthRequest,
			})
			c.Abort()
			return
		}

		hash := md5.Sum([]byte(payload.Password))
		md5Hex := strings.ToUpper(hex.EncodeToString(hash[:]))

		user, err := repositories.NewUserRepository().GetUserWS(c, db, payload.Username)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": constants.CodeErrorSendMidTier,
				"status_desc": constants.StatusInvalidAuthRequest,
			})
			c.Abort()
			return
		}

		if md5Hex != strings.ToUpper(user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": constants.CodeErrorSendMidTier,
				"status_desc": constants.StatusInvalidAuthRequest,
			})
			c.Abort()
			return
		}

		c.Set("username", user.Username)

		c.Next()
	}
}
