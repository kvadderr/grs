package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string     `yaml:"env"`
	GRPC   GRPCConfig `yaml:"grpc"`
	HTTP HTTPConfig `yaml:"http"`
	Database DatabaseConfig `yaml:"database"`
	Secret string     `yaml:"secret"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port string `yaml:"port" env-default:"5432"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
}

func Load() *Config {
	configPath := getConfigPath()

	if configPath == "" {
		panic("provide config path")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist. Provided: " + configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic("wrong config file. Provided: " + configPath)
	}

	return &config
}

func getConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	return path
}