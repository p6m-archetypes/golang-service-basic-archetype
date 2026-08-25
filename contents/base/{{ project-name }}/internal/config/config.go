package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host           string
	Port           int
	ManagementPort int
	LoggingJSON    bool
	ServiceName    string
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "{{ service_port }}"))
	if err != nil {
		return nil, fmt.Errorf("config: PORT: %w", err)
	}
	mgmtPort, err := strconv.Atoi(getEnv("MANAGEMENT_PORT", "{{ management_port }}"))
	if err != nil {
		return nil, fmt.Errorf("config: MANAGEMENT_PORT: %w", err)
	}
	cfg := &Config{
		Host:           getEnv("HOST", "0.0.0.0"),
		Port:           port,
		ManagementPort: mgmtPort,
		LoggingJSON:    getEnv("LOGGING_STRUCTURED", "false") == "true",
		ServiceName:    getEnv("SERVICE_NAME", "{{ project-name }}"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
