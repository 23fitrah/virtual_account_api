package services

import (
	"fmt"
	"math/rand"
	"net/http"
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

type VirtualAccountService struct {
	db                       *gorm.DB
	redis                    *redis.Client
	virtualAccountRepository *repositories.VirtualAccountRepository
	cfg                      *config.ConfigVa
}

func NewVirtualAccountService(
	db *gorm.DB,
	redis *redis.Client,
	virtualAccountRepository *repositories.VirtualAccountRepository,
	cfg *config.ConfigVa,
) *VirtualAccountService {
	return &VirtualAccountService{
		db:                       db,
		redis:                    redis,
		virtualAccountRepository: virtualAccountRepository,
		cfg:                      cfg,
	}
}

func (s *VirtualAccountService) CreateVA(c *gin.Context, param *validations.CreateVAValidation) (resources.GeneralResponse[resources.CreateVAResource], int) {
	var errRepo error
	now := time.Now()
	expiredAt := now.Add(time.Duration(s.cfg.Va.ExpiredHours) * time.Hour)

	data := &model.VirtualAccount{
		ID:           uuid.NewString(),
		VANumber:     s.generateVANumber(),
		CustomerID:   param.RequestData.CustomerID,
		CustomerName: param.RequestData.CustomerName,
		Amount:       param.RequestData.Amount,
		Description:  param.RequestData.Description,
		Status:       1,
		ReferenceID:  param.RequestData.ReferenceID,
		ExpiredAt:    expiredAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	errRepo = s.virtualAccountRepository.DoCreateVA(c, data, s.db)
	if errRepo != nil {
		return resources.GeneralResponse[resources.CreateVAResource]{
			BaseResponse: resources.BaseResponse{
				Status:       constants.StatusCodeVaFailed,
				ResponseCode: constants.CodeVaFailed,
				Message:      constants.StatusErrorCustom + "Create Virtual Account",
				Errors:       errRepo.Error(),
			},
		}, http.StatusInternalServerError
	}
	vaData := resources.ToFormModelResource(data)
	return resources.GeneralResponse[resources.CreateVAResource]{
		BaseResponse: resources.BaseResponse{
			Status:       constants.StatusCodeVaSuccess,
			ResponseCode: constants.StatusCodeVaCreate,
			Message:      constants.StatusCreateVASuccess,
		},
		Data: &vaData,
	}, http.StatusOK
}

func (s *VirtualAccountService) GetVAStatus(c *gin.Context, vaNumber string) (resources.GeneralResponse[resources.GetVAResource], int) {
	var result *resources.GetVAResource
	var errRepo error

	result, errRepo = s.virtualAccountRepository.DoGetVAStatus(c, vaNumber, s.db)
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

func (s *VirtualAccountService) GetVA(c *gin.Context, custId string, status string, page, limit, offset int) (utils.PaginatedResponse, int) {
	var result []*resources.GetVAListResource
	var errRepo error
	var total int64

	result, total, errRepo = s.virtualAccountRepository.DoGetVA(c, custId, status, limit, offset, s.db)
	if errRepo != nil {
		return utils.PaginatedResponse{
			Status:       constants.StatusCodeVaFailed,
			ResponseCode: constants.CodeVaFailed,
			Message:      constants.StatusErrorCustom + "Get Virtual Account",
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

func (s *VirtualAccountService) generateVANumber() string {
	// Format: <prefix><timestamp><random4digit>
	ts := time.Now().Format("20060102150405")
	rnd := fmt.Sprintf("%04d", rand.Intn(10000))
	return s.cfg.Va.Prefix + ts + rnd
}
