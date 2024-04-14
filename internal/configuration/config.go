package configuration

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Network *NetworkCfg `yaml:"network"`
}

type NetworkCfg struct {
	ServerAddress string `yaml:"address"`
	BaseURL       string `yaml:"base-url"`
}

func LoadCfg(filePath string) (*Config, error) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error during configuration")
	}
	var cfg Config
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}
