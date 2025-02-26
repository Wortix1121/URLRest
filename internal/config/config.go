package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"env" env-default:"local"`
	StoragePaths  `yaml:"storages" `
	HTTPServer    `yaml:"http_server"`
	Authorization `yaml:"authorization"`
	AliasLength   int `yaml:"aliasLenght"`
}

type StoragePaths struct {
	StoragePath   string `yaml:"storage_path" env-required:"true"`
	StoragePathPg string `yaml:"storage_pathPG" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Authorization struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true" env:"HTTP_AUTORIZATION_PASSWORD"`
}

func MustLoad() *Config {
	var cfg Config

	err := cleanenv.ReadConfig("config/local.yaml", &cfg)
	if err != nil {
		log.Fatal("cannot read config: ", err)

	}

	return &cfg
}
