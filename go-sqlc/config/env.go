package config

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	values map[string]string
}

var (
	envConfig *EnvConfig
	once      sync.Once
)

func GetEnvConfig() *EnvConfig {
	once.Do(func() {
		envConfig = &EnvConfig{
			values: make(map[string]string),
		}

		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not found: %v", err)
		}

		envMap, err := godotenv.Read()
		if err == nil {
			envConfig.values = envMap
		}
	})
	return envConfig
}

func (e *EnvConfig) GetString(key string) (string, bool) {
	if value, exists := e.values[key]; exists && value != "" {
		return value, true
	}
	return "", false
}

func (e *EnvConfig) GetInt(key string) (int, bool) {
	strValue, ok := e.GetString(key)
	if !ok {
		return 0, false
	}
	if value, err := strconv.Atoi(strValue); err == nil {
		return value, true
	}
	return 0, false
}

func (e *EnvConfig) GetBool(key string) (bool, bool) {
	strValue, ok := e.GetString(key)
	if !ok {
		return false, false
	}
	if value, err := strconv.ParseBool(strValue); err == nil {
		return value, true
	}
	return false, false
}

func (e *EnvConfig) GetDuration(key string) (time.Duration, bool) {
	strValue, ok := e.GetString(key)
	if !ok {
		return 0, false
	}
	if value, err := time.ParseDuration(strValue); err == nil {
		return value, true
	}
	return 0, false
}

func (e *EnvConfig) GetValues() map[string]string {
	return e.values
}
