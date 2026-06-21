//go:build wireinject

package providers

import (
	"virtual_account_api/internal/handlers"

	"github.com/google/wire"
)

var HandlerProviderSet = wire.NewSet(
	handlers.NewVirtualAccountHandler,
	handlers.NewPaymentHandler,
)
