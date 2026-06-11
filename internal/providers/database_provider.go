//go:build wireinject

package providers

import (
	"virtual_account_api/config"
	"time"

	"github.com/google/wire"
)

func ProvideESConfig(cfg *config.Config) config.ESConfig {
	return config.ESConfig{
		URL:      cfg.Logging.ElasticSearch.URL,
		Username: cfg.Logging.ElasticSearch.Username,
		Password: cfg.Logging.ElasticSearch.Password,
		Timeout:  10 * time.Second,
	}
}

var DatabaseProviderSet = wire.NewSet(
	config.ConnectGormDB,
	config.ConnectRedis,
	config.Load,
	ProvideESConfig,
	config.NewElasticsearchConnection,
)
