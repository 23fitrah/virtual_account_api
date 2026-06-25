package middleware

import (
	"time"
	"virtual_account_api/constants"
	"virtual_account_api/models"
	"virtual_account_api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		oldVal, _ := c.Get(constants.ContextKeyOldValue)
		newVal, _ := c.Get(constants.ContextKeyNewValue)
		respMsg, _ := c.Get(constants.ContextKeyResponseMessage)
		username, _ := c.Get(constants.ContextKeyUser)
		menu, _ := c.Get(constants.ContextKeyMenu)

		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		var userStr string
		if username != nil {
			if u, ok := username.(string); ok {
				userStr = u
			}
		}

		go func(o, n, r, mn interface{}, u, p, m, ip string) {
			oldStr := utils.ToString(o)
			newStr := utils.ToString(n)
			rspStr := utils.ToString(r)
			menuStr := utils.ToString(mn)

			historyLog := models.HistoryLog{
				Timestamp:       time.Now().Local(),
				User:            u,
				Menu:            menuStr,
				Action:          m,
				NewValue:        newStr,
				OldValue:        oldStr,
				ResponseMessage: rspStr,
				IpAddress:       ip,
			}

			err := utils.InsertHistoryLog(c, db, historyLog)
			if err != nil {
				utils.LogError(constants.StatusFailedInsertHistoryLog, err)
			}

		}(oldVal, newVal, respMsg, menu, userStr, path, method, clientIP)
	}
}
