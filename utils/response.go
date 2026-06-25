package utils

import (
	"virtual_account_api/constants"

	"github.com/gin-gonic/gin"
)

func Responds[T any](c *gin.Context, resp T, status int) {
	c.Set(constants.ContextKeyResponseMessage, resp)
	c.JSON(status, resp)
}
