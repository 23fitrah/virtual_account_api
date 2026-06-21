//go:build wireinject

package providers

import (
	"virtual_account_api/internal/repositories"

	"github.com/google/wire"
)

var RepositoryProviderSet = wire.NewSet(
	repositories.NewVirtualAccountRepository,
	repositories.NewPaymentRepository,
)
