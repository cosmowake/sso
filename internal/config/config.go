package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default: "local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-default: "1h"`
	GRPC        GRPC          `yaml:"grpc"`
	Mongo       Mongo         `yaml:"mongo"`
}

type GRPC struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout" env-default: "1h"`
}

type Mongo struct {
	Address string `yaml: "address" env-default:"mongodb://localhost:27017"`
}

func MustLoad() *Config {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.Parse()

	if configPath == "" {
		panic("config path error")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file is not exist error")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config file")
	}

	return &cfg
}
