//go:build wireinject

package providers

import (
	"virtual_account_api/internal/services"

	"github.com/google/wire"
)

var ServiceProviderSet = wire.NewSet(
	services.NewVirtualAccountService,
	services.NewPaymentService,
)
