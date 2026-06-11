package config

import (
	"os"
	"strconv"
)

type ConfigVa struct {
	Va VAConfig
}

type VAConfig struct {
	Prefix       string
	ExpiredHours int
}

func LoadVA() (*ConfigVa, error) {
	vaExpired, _ := strconv.Atoi(os.Getenv("VA_EXPIRED_HOURS"))

	return &ConfigVa{
		Va: VAConfig{
			Prefix:       os.Getenv("VA_PREFIX"),
			ExpiredHours: vaExpired,
		},
	}, nil
}
