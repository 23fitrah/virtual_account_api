package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"virtual_account_api/constants"
	"virtual_account_api/internal/repositories"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthWSMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": constants.CodeErrorSendMidTier,
				"statusDesc": constants.StatusInvalidRequestBody,
			})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := sonic.Unmarshal(bodyBytes, &payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": constants.CodeErrorSendMidTier,
				"statusDesc": constants.StatusInvalidAuthRequest,
			})
			c.Abort()
			return
		}

		if payload.Username == "" || payload.Password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": constants.CodeErrorSendMidTier,
				"statusDesc": constants.StatusInvalidAuthRequest,
			})
			c.Abort()
			return
		}

		user, err := repositories.NewUserRepository().GetUserWS(c, db, payload.Username)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": constants.CodeErrorSendMidTier,
				"statusDesc": constants.StatusInvalidAuthRequest,
			})
			c.Abort()
			return
		}

		hash := md5.Sum([]byte(payload.Password))
		md5Hex := strings.ToUpper(hex.EncodeToString(hash[:]))

		if md5Hex != strings.ToUpper(user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": constants.CodeErrorSendMidTier,
				"statusDesc": constants.StatusInvalidAuthRequest,
			})
			c.Abort()
			return
		}

		c.Set("username", user.Username)

		c.Next()
	}
}
