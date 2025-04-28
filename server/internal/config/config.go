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
	App struct {
		Port int `mapstructure:"Port"`
	} `mapstructure:"app"`

	Database struct {
		Source string `mapstructure:"Source"`

		SQLite struct {
			Path string `mapstructure:"Path"`
		} `mapstructure:"SQLite"`

		MySQL struct {
			Host            string `mapstructure:"Host"`
			Port            int    `mapstructure:"Port"`
			Name            string `mapstructure:"Name"`
			Username        string `mapstructure:"Username"`
			Password        string `mapstructure:"Password"`
			MaxIdleConns    int    `mapstructure:"MaxIdleConns"`
			SetMaxOpenConns int    `mapstructure:"SetMaxOpenConns"`
		} `mapstructure:"MySQL"`
	} `mapstructure:"database"`

	Worker struct {
		Endpoints []string `mapstructure:"Endpoints"`
		Cpu       bool     `mapstructure:"Cpu"`
	} `mapstructure:"worker"`

	Minio struct {
		Endpoint      string `mapstructure:"Endpoint"`
		AccessKey     string `mapstructure:"AccessKey"`
		SecretKey     string `mapstructure:"SecretKey"`
		DefaultBucket string `mapstructure:"DefaultBucket"`
	} `mapstructure:"minio"`
}
