package config

import (
	"fmt"

	"gorm.io/gorm"
)

type Config struct {
	Logging LoggingConfig `json:"logging"`
}

type LoggingConfig struct {
	ElasticSearch ElasticearchConfig `json:"elasticSearch"`
	Level         string             `json:"level"`
}

type ElasticearchConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}

func Load() (*Config, error) {
	esUrl := GetEnv("ES_URL_SA", "")
	esUser := GetEnv("ES_USER_SA", "")
	esPwd := GetEnv("ES_PASS_SA", "")
	esEnbl := GetEnv("ES_ENABLED_SA", "")
	esLvl := GetEnv("ES_LEVEL_SA", "")

	return &Config{
		Logging: LoggingConfig{
			ElasticSearch: ElasticearchConfig{
				URL:      esUrl,
				Username: esUser,
				Password: esPwd,
				Enabled:  esEnbl == "true",
			},
			Level: esLvl,
		},
	}, nil
}

func GetParameterByType(db *gorm.DB, jenis string) (string, error) {
	var status string

	err := db.Table("SA_PARAMETER").Select("VALUE").Where("NAME = ?", jenis).Scan(&status).Error

	if err != nil {
		return "", fmt.Errorf("failed to query W_PARAMETER for jenis=%s: %w", jenis, err)
	}

	if status == "" {
		return "", fmt.Errorf("type not found: %s", jenis)
	}

	return status, nil
}
