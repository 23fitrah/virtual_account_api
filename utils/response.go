package utils

import (
	"github.com/gin-gonic/gin"
)

func Responds[T any](c *gin.Context, resp T, status int) {
	c.JSON(status, resp)
}
