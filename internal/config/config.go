// Package config defines the config loading for the app
package config

import (
	"errors"
	"github.com/mfamador/grpc-input/internal/server"
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

// LoggerConfig has the log level and if we should pretty print
type loggerConfig struct {
	Level  string `yaml:"level"`
	Pretty bool   `yaml:"pretty"`
}

// AppConfig is main app config
type AppConfig struct {
	Logger loggerConfig  `yaml:"logger"`
	Server server.Config `yaml:"server"`
}

var (
	configFiles = []string{"config.yaml"}
	// Config contains all configuration values
	Config AppConfig
)

func searchConfig(dir string) (string, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	dirPath := filepath.Join(absPath, "config")
	if fileInfo, err := os.Stat(dirPath); err == nil && fileInfo.IsDir() {
		return dirPath, nil
	}

	if absPath == "/" {
		return "", errors.New("not found")
	}

	return searchConfig(filepath.Join(absPath, ".."))
}

// NewConfig returns a new config
func init() {
	setEnvVars()
	configDir := os.Getenv("CONFIGOR_DIR")

	if configDir == "" {
		var err error
		if configDir, err = searchConfig("."); err != nil {
			panic("Config dir not found")
		}
	}

	for i, v := range configFiles {
		configFiles[i] = filepath.Join(configDir, v)
	}

	config := configor.New(&configor.Config{ENVPrefix: "-"})
	if err := config.Load(&Config, configFiles...); err != nil {
		panic("Invalid config")
	}
}

func setEnvVars() {
	if os.Getenv("DB_READINGSSTORAGECONNECTIONSTRING") == "" {
		os.Setenv("DB_READINGSSTORAGECONNECTIONSTRING", os.Getenv("INTEGRATION_TEST_STORAGECONNECTIONSTRING"))     //nolint
		os.Setenv("DB_PACKETSSTORAGECONNECTIONSTRING", os.Getenv("INTEGRATION_TEST_STORAGECONNECTIONSTRING"))      //nolint
		os.Setenv("DB_EVENTSSTORAGECONNECTIONSTRING", os.Getenv("INTEGRATION_TEST_STORAGECONNECTIONSTRING"))       //nolint
		os.Setenv("DB_CONSUMPTIONSSTORAGECONNECTIONSTRING", os.Getenv("INTEGRATION_TEST_STORAGECONNECTIONSTRING")) //nolint
		os.Setenv("DB_REFILLSSTORAGECONNECTIONSTRING", os.Getenv("INTEGRATION_TEST_STORAGECONNECTIONSTRING"))      //nolint
	}
}
