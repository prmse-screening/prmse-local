package config

import (
	"github.com/spf13/viper"
)

var Cfg *Config

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic(err)
	}
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Worker   WorkerConfig
	Minio    MinioConfig
}

type AppConfig struct {
	Port int
}

type DatabaseConfig struct {
	Path string
}

type WorkerConfig struct {
	Endpoints []string
	Cpu       bool
}

type MinioConfig struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	DefaultBucket string
}
