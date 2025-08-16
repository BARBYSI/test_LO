package config

import (
	"encoding/json"
	"fmt"
	"os"

	"task-manager/pkg/logger"
)

type Config struct {
	AppPort       string `json:"app_port"`
	LoggerEnabled bool   `json:"logger_enabled"`
	// feel free to add more fields
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("../config.json")
	if err != nil {
		logger.LogError(fmt.Sprintf("failed to load config: %v", err))
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		logger.LogError(fmt.Sprintf("failed to unmarshal config: %v", err))
		return nil, err
	}

	return &config, nil
}
