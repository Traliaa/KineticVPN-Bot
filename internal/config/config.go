package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	configFilePathENV = "CONFIG_FILE"
	tokenTelegramENV  = "TELEGRAM_TOKEN"
)

type config struct {
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	//Jaeger struct {
	//	Host string `yaml:"host"`
	//	Port int    `yaml:"port"`
	//} `yaml:"jaeger"`
	//ProductService struct {
	//	Host  string `yaml:"host"`
	//	Port  int    `yaml:"port"`
	//	Token string `yaml:"token"`
	//	Limit int    `yaml:"limit"`
	//	Burst int    `yaml:"burst"`
	//} `yaml:"product_service"`
	//LomsService struct {
	//	Host string `yaml:"host"`
	//	Port int    `yaml:"port"`
	//} `yaml:"loms_service"`
	//Redis struct {
	//	Host string `yaml:"host"`
	//	Port int    `yaml:"port"`
	//} `yaml:"redis"`
}

// NewNoop returns empty struct.
func newNoop() config {
	return config{}
}

func NewConfig() config {
	file, err := os.Open(os.Getenv(configFilePathENV))
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	config := config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	token := os.Getenv(tokenTelegramENV)
	if token != "" {
		config.Telegram.Token = token
	}

	return config
}
