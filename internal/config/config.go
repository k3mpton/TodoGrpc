package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `yaml:"env" env-default:"local"`
	Grpc gRpc   `yaml:"gRPC" env-required:"true"`
}

type gRpc struct {
	Port    int           `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustConfigLoad() Config {
	path := GetPath()

	if _, err := os.Stat(path); err != nil {
		panic("не удалось найти файл по заданному пути")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("fail read config")
	}
	return cfg
}

var (
	PathConfig = flag.String("cfg", "./local/config.yaml", "config path")
)

func GetPath() (path string) {
	flag.Parse()
	return *PathConfig
}
