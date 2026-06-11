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

func (s *VirtualAccountService) CreateVA(c *gin.Context, param *validations.CreateVirtualAccountValidation) (resources.GeneralResponse[resources.CreateVAResource], int) {
	var result resources.CreateVAResource
	var errRepo error
	now := time.Now()
	expiredAt := now.Add(time.Duration(s.cfg.Va.ExpiredHours) * time.Hour)

	data := &model.VirtualAccount{
		ID:           uuid.NewString(),
		VANumber:     s.generateVANumber(),
		CustomerID:   param.CustomerID,
		CustomerName: param.CustomerName,
		Amount:       param.Amount,
		Description:  param.Description,
		Status:       constants.VaStatusPending,
		ReferenceID:  param.ReferenceID,
		ExpiredAt:    expiredAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	errRepo = s.virtualAccountRepository.DoCreateVA(c, data, s.db)
	if errRepo != nil {
		return resources.GeneralResponse[resources.CreateVAResource]{
			BaseResponse: resources.BaseResponse{
				StatusCode: constants.CodeErrorSendMidTier,
				StatusDesc: constants.StatusErrorCustom + errRepo.Error(),
			},
		}, http.StatusInternalServerError
	}

	return resources.GeneralResponse[resources.CreateVAResource]{
		BaseResponse: resources.BaseResponse{
			StatusCode: constants.CodeVaCreate,
			StatusDesc: constants.StatusGetSuccess,
		},
		Data: result,
	}, http.StatusOK
}

func (s *VirtualAccountService) GetVAStatus(c *gin.Context, vaNumber string) (resources.GeneralResponse[resources.GetVAResource], int) {
	var result *resources.GetVAResource
	var errRepo error

	result, errRepo = s.virtualAccountRepository.DoGetVAStatus(c, vaNumber, s.db)
	if errRepo != nil {
		return resources.GeneralResponse[resources.GetVAResource]{
			BaseResponse: resources.BaseResponse{
				StatusCode: constants.CodeErrorSendMidTier,
				StatusDesc: constants.StatusErrorCustom + errRepo.Error(),
			},
		}, http.StatusInternalServerError
	}
	return resources.GeneralResponse[resources.GetVAResource]{
		BaseResponse: resources.BaseResponse{
			StatusCode: constants.CodeVaGetStatus,
			StatusDesc: constants.StatusGetSuccess,
		},
		Data: *result,
	}, http.StatusOK
}

func (s *VirtualAccountService) GetVA(c *gin.Context, custId string, status string, params utils.PaginationParams) (utils.PaginationResult, int) {
	var result resources.GetVAListResource
	var errRepo error
	counterpart := params.Filter["counterpart"]
	startDate := params.Filter["startDate"]
	endDate := params.Filter["endDate"]

	result, errRepo = s.virtualAccountRepository.DoGetVA(c, s.db, custId, status, params)
	if errRepo != nil {
		return resources.GeneralResponse[resources.GetVAListResource]{
			BaseResponse: resources.BaseResponse{
				StatusCode: constants.CodeErrorSendMidTier,
				StatusDesc: constants.StatusErrorCustom + errRepo.Error(),
			},
		}, http.StatusInternalServerError
	}
	return resources.GeneralResponse[resources.GetVAListResource]{
		BaseResponse: resources.BaseResponse{
			StatusCode: constants.CodeVaGet,
			StatusDesc: constants.StatusGetSuccess,
		},
		Data: result,
	}, http.StatusOK
}

func (s *VirtualAccountService) generateVANumber() string {
	// Format: <prefix><timestamp><random4digit>
	ts := time.Now().Format("20060102150405")
	rnd := fmt.Sprintf("%04d", rand.Intn(10000))
	return s.cfg.Va.Prefix + ts + rnd
}
