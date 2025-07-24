package config

import (
	"os"
	"strconv"
)

type Config struct {
	VolcengineAPIKey  string
	VolcengineRegion  string
	VolcengineBaseURL string
	Port              int
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, err
	}

	return &Config{
		VolcengineAPIKey:  getEnv("VOLCENGINE_API_KEY", "dc083843-06ad-4553-a7ce-4e86a06c4235"),
		VolcengineRegion:  getEnv("VOLCENGINE_REGION", "cn-beijing"),
		VolcengineBaseURL: getEnv("VOLCENGINE_BASE_URL", "https://ark.cn-beijing.volces.com/api/v3"),
		Port:              port,
	}, nil

}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
