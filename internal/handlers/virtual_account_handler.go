package handlers

import (
	"virtual_account_api/internal/services"
	"virtual_account_api/internal/validations"
	"virtual_account_api/utils"

	"github.com/gin-gonic/gin"
)

type VirtualAccountHandler struct {
	virtualAccountService *services.VirtualAccountService
}

func NewVirtualAccountHandler(
	virtualAccountService *services.VirtualAccountService,
) *VirtualAccountHandler {
	return &VirtualAccountHandler{
		virtualAccountService: virtualAccountService,
	}
}

func (h *VirtualAccountHandler) CreateVA(c *gin.Context) {
	req, respErr, code := utils.ValidateAndBind[validations.CreateVAValidation](c)
	if respErr != nil {
		utils.Responds(c, respErr, code)
		return
	}

	response, respCode := h.virtualAccountService.CreateVA(c, req)

	utils.Responds(c, response, respCode)
}

func (h *VirtualAccountHandler) GetVAStatus(c *gin.Context) {
	vaNumber := c.Param("va_number")
	response, respCode := h.virtualAccountService.GetVAStatus(c, vaNumber)

	utils.Responds(c, response, respCode)
}

func (h *VirtualAccountHandler) GetVA(c *gin.Context) {
	custId := c.Query("customer_id")
	status := c.Query("status")
	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	response, respCode := h.virtualAccountService.GetVA(c, custId, status, page, limit, offset)

	utils.Responds(c, response, respCode)
}
