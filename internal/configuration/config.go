// Package configuration provides functionality to load and retrieve
// configuration data for the short-url application.
package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (

	// ServerAddress is a constant representing the key for the server address configuration in the configuration file.
	ServerAddress = "server_address"

	// BaseURL is a constant representing the key for the base URL configuration in the configuration file.
	BaseURL = "base_url"

	// FileStorePath is a constant representing the key for the file storage path configuration in the configuration file.
	FileStorePath = "file_storage_path"

	// DataBaseDNS is a constant representing the key for the database DNS configuration in the configuration file.
	DataBaseDNS = "database_dsn"

	// EnableHTTPS = "enable_https"
	EnableHTTPS = "enable_https"
)

// Config holds all the configuration data needed for the
// short-url application to run correctly.
type Config struct {
	Network *NetworkCfg `yaml:"network"`
	Logger  *LoggerCfg  `yaml:"logger"`
	Storage *Storage    `yaml:"storage"`
}

// LoggerCfg structure
type LoggerCfg struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

// NetworkCfg structure
type NetworkCfg struct {
	ServerAddress string `yaml:"address"`
	BaseURL       string `yaml:"base-url"`
	EnableHTTPS   bool   `yaml:"enable-https"`
}

// Storage structure
type Storage struct {
	FileStoragePath string `yaml:"address"`
	ConnString      string
	PGPoolCfg       *pgxpool.Config
}

// New creates a new Config instance, loading data from
// environment variables or configuration files as needed.
// Returns an error if required configuration data could not be loaded.
func New() *Config {
	networkCfg := NetworkCfg{}
	loggerCfg := LoggerCfg{}
	storageCfg := Storage{}
	var filePath string

	flag.StringVar(&networkCfg.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&networkCfg.BaseURL, "b", "http://localhost:8080", "base address")
	flag.StringVar(&loggerCfg.Level, "l", "info", "logger level")
	flag.StringVar(&storageCfg.FileStoragePath, "f", "/tmp/short-url-db.json", "file to store data")
	flag.StringVar(&storageCfg.ConnString, "d", "", "file to store data")
	flag.StringVar(&filePath, "c", "", "path to config file")
	flag.Parse()

	if filePath == "" {
		if filePathEnv := os.Getenv("CONFIG"); filePathEnv != "" {
			filePath = filePathEnv
		}
	}

	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		networkCfg.ServerAddress = serverAddress
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		networkCfg.BaseURL = baseURL
	}
	if loggerLevel := os.Getenv("LOG_LEVEL"); loggerLevel != "" {
		loggerCfg.Level = loggerLevel
	}

	if filaStoragePath := os.Getenv("FILE_STORAGE_PATH"); filaStoragePath != "" {
		storageCfg.FileStoragePath = filaStoragePath
	}
	if connString := os.Getenv("DATABASE_DSN"); connString != "" {
		storageCfg.ConnString = connString
	}
	cfg := &Config{
		Network: &networkCfg,
		Logger:  &loggerCfg,
		Storage: &storageCfg,
	}
	readFromFile(filePath, cfg)
	return cfg
}

// Handle error
func readFromFile(filePath string, config *Config) {
	if filePath != "" {
		f, err := os.ReadFile(filePath)
		if err != nil {
			return
		}
		var data map[string]interface{}
		err = json.Unmarshal(f, &data)
		if err != nil {
			return
		}
		for k, v := range data {
			switch k {
			case ServerAddress:
				if config.Network.ServerAddress == "" {
					config.Network.ServerAddress = v.(string)
				}

			case BaseURL:
				if config.Network.BaseURL == "" {
					config.Network.BaseURL = v.(string)
				}
			case FileStorePath:
				if config.Storage.FileStoragePath == "" {
					config.Storage.FileStoragePath = v.(string)
				}
			case DataBaseDNS:
				if config.Storage.ConnString == "" {
					config.Storage.ConnString = v.(string)
				}

			case EnableHTTPS:
				val := v.(bool)
				if !config.Network.EnableHTTPS && val {
					config.Network.EnableHTTPS = val
				}
			default:
				fmt.Printf("field of config is unknown%T!\n", v)
			}
		}
	}
}
