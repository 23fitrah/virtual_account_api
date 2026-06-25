package routes

import (
	"virtual_account_api/internal/injector"
	"virtual_account_api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterV1Router(rg *gin.RouterGroup, ct *injector.AppContainer) {
	v1 := rg.Group("/v1")

	v1.Use(middleware.AuthWSMiddleware(ct.DB), middleware.AuditMiddleware(ct.DB))

	v1.POST("/virtual-accounts/create", ct.VirtualAccountHandler.CreateVA)

	v1.POST("/virtual-accounts/:va_number", ct.VirtualAccountHandler.GetVAStatus)

	v1.POST("/virtual-accounts", ct.VirtualAccountHandler.GetVA)

	v1.POST("/payments", ct.PaymentHandler.CallbackPayment)

	v1.POST("/payments/history", ct.PaymentHandler.GetPaymentHistory)

}
