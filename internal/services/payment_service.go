package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"virtual_account_api/config"
	"virtual_account_api/constants"
	"virtual_account_api/internal/repositories"
	"virtual_account_api/internal/validations"
	model "virtual_account_api/models"
	"virtual_account_api/resources"
	"virtual_account_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PaymentService struct {
	db                *gorm.DB
	redis             *redis.Client
	paymentRepository *repositories.PaymentRepository
	cfg               *config.ConfigVa
}

func NewPaymentService(
	db *gorm.DB,
	redis *redis.Client,
	paymentRepository *repositories.PaymentRepository,
	cfg *config.ConfigVa,
) *PaymentService {
	return &PaymentService{
		db:                db,
		redis:             redis,
		paymentRepository: paymentRepository,
		cfg:               cfg,
	}
}

func (s *PaymentService) CallbackPayment(c *gin.Context, param *validations.CallbackPaymentValidation) (resources.GeneralResponse[resources.CreatePaymentResource], int) {
	var errRepo error
	now := time.Now()

	va, err := s.paymentRepository.DoGetVAStatus(c, param.RequestData.VANumber, s.db)
	if err != nil {
		return resources.GeneralResponse[resources.CreatePaymentResource]{
			BaseResponse: resources.BaseResponse{
				Status:       constants.StatusCodePaymentFailed,
				ResponseCode: constants.CodePaymentFailed,
				Message:      constants.StatusErrorCustom + "Get Virtual Account Status",
				Errors:       err.Error(),
			},
		}, http.StatusInternalServerError
	}

	if va.StatusName == constants.VaStatusPaid {
		return resources.GeneralResponse[resources.CreatePaymentResource]{
			BaseResponse: resources.BaseResponse{
				Status:       constants.StatusCodePaymentFailed,
				ResponseCode: constants.CodePaymentFailed,
				Message:      "Virtual account already paid",
			},
		}, http.StatusBadRequest
	}

	vaExpiredAt, errParse := time.Parse(time.RFC3339, va.ExpiredAt)
	if errParse != nil {
		return resources.GeneralResponse[resources.CreatePaymentResource]{
			BaseResponse: resources.BaseResponse{
				Status:       constants.StatusCodePaymentFailed,
				ResponseCode: constants.CodePaymentFailed,
				Message:      constants.StatusErrorCustom + " Virtual Account Expired",
				Errors:       errParse.Error(),
			},
		}, http.StatusInternalServerError
	}

	if va.StatusName == constants.VaStatusExpired || time.Now().After(vaExpiredAt) {
		return resources.GeneralResponse[resources.CreatePaymentResource]{
			BaseResponse: resources.BaseResponse{
				Status:       constants.StatusCodePaymentFailed,
				ResponseCode: constants.CodePaymentFailed,
				Message:      "Virtual account has expired",
			},
		}, http.StatusBadRequest
	}
	amountVa, errParseAmount := strconv.ParseFloat(va.Amount, 64)
	if errParseAmount != nil {
		return resources.GeneralResponse[resources.CreatePaymentResource]{
			BaseResponse: resources.BaseResponse{
				Status:       constants.StatusCodePaymentFailed,
				ResponseCode: constants.CodePaymentFailed,
				Message:      constants.StatusErrorCustom + "Match Amount Virtual Account",
				Errors:       errParseAmount.Error(),
			},
		}, http.StatusInternalServerError
	}
	if param.RequestData.PaidAmount != amountVa {
		return resources.GeneralResponse[resources.CreatePaymentResource]{
			BaseResponse: resources.BaseResponse{
				Status:       constants.StatusCodePaymentFailed,
				ResponseCode: constants.CodePaymentFailed,
				Message:      "Paid amount does not match virtual account amount"},
		}, http.StatusBadRequest
	}

	rawPayload, _ := json.Marshal(param.RequestData)
	data := &model.Payment{
		ID:               uuid.NewString(),
		VANumber:         param.RequestData.VANumber,
		VirtualAccountID: va.Id,
		PaidAmount:       param.RequestData.PaidAmount,
		PaymentChannel:   param.RequestData.PaymentChannel,
		Status:           2, // Assuming 2 means paid
		RawPayload:       string(rawPayload),
		CreatedAt:        now,
	}

	errRepo = s.paymentRepository.DoPaymentCallback(c, data, s.db)
	if errRepo != nil {
		return resources.GeneralResponse[resources.CreatePaymentResource]{
			BaseResponse: resources.BaseResponse{
				Status:       constants.StatusCodePaymentFailed,
				ResponseCode: constants.CodePaymentFailed,
				Message:      constants.StatusErrorCustom + "Create Payment",
				Errors:       errRepo.Error(),
			},
		}, http.StatusInternalServerError
	}
	paymentData := resources.ToFormModelPaymentResource(data)
	return resources.GeneralResponse[resources.CreatePaymentResource]{
		BaseResponse: resources.BaseResponse{
			Status:       constants.StatusCodePaymentSuccess,
			ResponseCode: constants.CodePaymentSuccess,
			Message:      constants.StatusPaymentSuccess,
		},
		Data: &paymentData,
	}, http.StatusOK
}

func (s *PaymentService) GetVAStatus(c *gin.Context, vaNumber string) (resources.GeneralResponse[resources.GetVAResource], int) {
	var result *resources.GetVAResource
	var errRepo error

	result, errRepo = s.paymentRepository.DoGetVAStatus(c, vaNumber, s.db)
	if errRepo != nil {
		return resources.GeneralResponse[resources.GetVAResource]{
			BaseResponse: resources.BaseResponse{
				Status:       constants.StatusCodeVaFailed,
				ResponseCode: constants.CodeVaFailed,
				Message:      constants.StatusErrorCustom + "Get Virtual Account Status",
				Errors:       errRepo.Error(),
			},
		}, http.StatusInternalServerError
	}
	return resources.GeneralResponse[resources.GetVAResource]{
		BaseResponse: resources.BaseResponse{
			Status:       constants.StatusCodeVaSuccess,
			ResponseCode: constants.CodeVaSuccess,
			Message:      constants.StatusGetSuccess,
		},
		Data: result,
	}, http.StatusOK
}

func (s *PaymentService) GetPaymentHistory(c *gin.Context, vaNumber string, page, limit, offset int) (utils.PaginatedResponse, int) {
	var result []*resources.GetPaymentListResource
	var errRepo error
	var total int64

	result, total, errRepo = s.paymentRepository.DoGetPaymentHistory(c, vaNumber, limit, offset, s.db)
	if errRepo != nil {
		return utils.PaginatedResponse{
			Status:       constants.StatusCodePaymentFailed,
			ResponseCode: constants.CodePaymentFailed,
			Message:      constants.StatusErrorCustom + "Get Payment History",
			Errors:       errRepo.Error(),
			Data:         result,
		}, http.StatusInternalServerError
	}

	pagination := utils.BuildPagination(page, limit, total)

	return utils.PaginatedResponse{
		Status:       constants.StatusCodeVaSuccess,
		ResponseCode: constants.CodeVaSuccess,
		Message:      constants.StatusGetSuccess,
		Data:         result,
		Pagination:   pagination,
	}, http.StatusOK
}

func (s *PaymentService) generateVANumber() string {
	// Format: <prefix><timestamp><random4digit>
	ts := time.Now().Format("20060102150405")
	rnd := fmt.Sprintf("%04d", rand.Intn(10000))
	return s.cfg.Va.Prefix + ts + rnd
}
