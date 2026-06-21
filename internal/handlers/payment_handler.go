package handlers

import (
	"virtual_account_api/internal/services"
	"virtual_account_api/internal/validations"
	"virtual_account_api/utils"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(
	paymentService *services.PaymentService,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CallbackPayment(c *gin.Context) {
	req, respErr, code := utils.ValidateAndBind[validations.CallbackPaymentValidation](c)
	if respErr != nil {
		utils.Responds(c, respErr, code)
		return
	}

	response, respCode := h.paymentService.CallbackPayment(c, req)

	utils.Responds(c, response, respCode)
}

func (h *PaymentHandler) GetVAStatus(c *gin.Context) {
	vaNumber := c.Param("va_number")

	_, respErr, code := utils.ValidateAndBind[validations.GetVAStatusValidation](c)
	if respErr != nil {
		utils.Responds(c, respErr, code)
		return
	}

	response, respCode := h.paymentService.GetVAStatus(c, vaNumber)

	utils.Responds(c, response, respCode)
}

func (h *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	vaNumber := c.Query("va_number")
	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	_, respErr, code := utils.ValidateAndBind[validations.GetPaymentHistoryValidation](c)
	if respErr != nil {
		utils.Responds(c, respErr, code)
		return
	}
	response, respCode := h.paymentService.GetPaymentHistory(c, vaNumber, page, limit, offset)

	utils.Responds(c, response, respCode)
}
