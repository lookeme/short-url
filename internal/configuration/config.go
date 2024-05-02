package configuration

import (
	"flag"
	"os"
)

type Config struct {
	Network *NetworkCfg `yaml:"network"`
	Logger  *LoggerCfg  `yaml:"logger"`
	Storage *Storage    `yaml:"storage"`
}
type LoggerCfg struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}
type NetworkCfg struct {
	ServerAddress string `yaml:"address"`
	BaseURL       string `yaml:"base-url"`
}

type Storage struct {
	FileStoragePath string `yaml:"address"`
}

func CreateConfig() *Config {
	networkCfg := NetworkCfg{}
	loggerCfg := LoggerCfg{}
	storageCfg := Storage{}
	flag.StringVar(&networkCfg.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&networkCfg.BaseURL, "b", "http://localhost:8080", "base address")
	flag.StringVar(&loggerCfg.Level, "l", "info", "logger level")
	flag.StringVar(&storageCfg.FileStoragePath, "f", "/tmp/short-url-db.json", "file to store data")

	flag.Parse()
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
	return &Config{
		Network: &networkCfg,
		Logger:  &loggerCfg,
		Storage: &storageCfg,
	}
}
