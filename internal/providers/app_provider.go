//go:build wireinject

package providers

import (
	"github.com/google/wire"
)

var AppProviderSet = wire.NewSet(
	DatabaseProviderSet,
	RepositoryProviderSet,
	ServiceProviderSet,
	HandlerProviderSet,
)
