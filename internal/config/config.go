package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	configFilePathENV = "CONFIG_FILE"
	tokenTelegramENV  = "TELEGRAM_TOKEN"
	databaseDSN       = "DATABASE_DSN"
)

// Config ...
type Config struct {
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	DB      string `yaml:"db_dsn"`
	Service struct {
		Host     string `yaml:"host"`
		HTTPPort int    `yaml:"http_port"`
	} `yaml:"service"`
}

// NewNoop returns empty struct.
func newNoop() Config {
	return Config{}
}

func NewConfig() *Config {

	configFileName := os.Getenv(configFilePathENV)
	if configFileName == "" {
		configFileName = "values_local.yaml"
	}
	file, err := os.Open("configs/" + configFileName)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	token := os.Getenv(tokenTelegramENV)
	if token != "" {
		config.Telegram.Token = token
	}

	dsn := os.Getenv(databaseDSN)
	if dsn != "" {
		config.DB = dsn
	}

	return &config
}
