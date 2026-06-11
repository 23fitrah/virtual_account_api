package routes

import (
	"virtual_account_api/internal/injector"
	"virtual_account_api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterV1Router(rg *gin.RouterGroup, ct *injector.AppContainer) {
	v1 := rg.Group("/v1")

	v1.Use(middleware.AuthWSMiddleware(ct.DB), middleware.AuditMiddleware(ct.DB))

	v1.POST("/virtual-accounts", ct.VirtualAccountHandler.CreateVA)

	v1.GET("/virtual-accounts/:va_number", ct.VirtualAccountHandler.GetVAStatus)

	v1.GET("/virtual-accounts", ct.VirtualAccountHandler.GetVA)

	/*	payment := v1.Group("/payments")
		{
			payment.POST("/process", ct.DataMt199Handler.Detail)
			payment.POST("/history", ct.DataMt199Handler.ClosePending)
		}*/
}
