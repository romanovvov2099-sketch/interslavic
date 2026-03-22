package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP        HTTP        `yaml:"http"`
	Logging     Logging     `yaml:"logging"`
	Database    Database    `yaml:"database"`
	JWT         JWT         `yaml:"jwt"`
}

type HTTP struct {
	Port             string `yaml:"port"`
	AuthorizationKey string `yaml:"authorization"`
}

type Logging struct {
	Dir        string `yaml:"dir"`
	Enable     bool   `yaml:"enable"`
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	SavingDays int    `yaml:"saving_days"`
}

type Database struct {
	URL string `yaml:"url"`
}

type JWT struct {
	SecretKey string `yaml:"secret_key"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig("./config.yaml", cfg); err != nil {
		fmt.Printf("failed get config: %s\n", err)
	}

	return cfg, nil
}