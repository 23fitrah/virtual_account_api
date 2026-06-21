//go:build wireinject

//go:generate wire

package injector

import (
	"virtual_account_api/internal/handlers"
	"virtual_account_api/internal/providers"
	"virtual_account_api/utils"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppContainer struct {
	DB                    *gorm.DB
	RedisClient           *redis.Client
	ValidationUtils       *utils.ValidationUtils
	VirtualAccountHandler *handlers.VirtualAccountHandler
	PaymentHandler        *handlers.PaymentHandler
}

// NewAppContainer creates an AppContainer with all dependencies
func NewAppContainer(
	db *gorm.DB,
	redisClient *redis.Client,
	validationUtils *utils.ValidationUtils,
	virtualAccountHandler *handlers.VirtualAccountHandler,
	paymentHandler *handlers.PaymentHandler,
) *AppContainer {
	return &AppContainer{
		DB:                    db,
		RedisClient:           redisClient,
		ValidationUtils:       validationUtils,
		PaymentHandler:        paymentHandler,
		VirtualAccountHandler: virtualAccountHandler,
	}
}

// InitializeApp creates the complete application with all dependencies
func InitializeApp() (*AppContainer, error) {
	wire.Build(
		providers.AppProviderSet,
		NewAppContainer,
		utils.NewValidationUtils,
	)
	return nil, nil
}
