package config

import "os"

type ConfigRepositoryContract interface {
	GetString(key string, defaultValue string) string
}

type ConfigRepository struct {
}

func (repo *ConfigRepository) GetString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
