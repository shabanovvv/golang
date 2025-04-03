package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf  `mapstructure:"logger"`
	Storage StorageConf `mapstructure:"storage"`
	Server  ServerConf  `mapstructure:"server"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
}

type StorageConf struct {
	// postgresDSN string `mapstructure:"postgres_dsn"`
	// timeout     int    `mapstructure:"timeout"`
}

type ServerConf struct {
	Host     string `mapstructure:"host"`
	GrpcPort int    `mapstructure:"grpc_port"`
	HTTPPort int    `mapstructure:"http_port"`
}

func NewConfig(fileName string) Config {
	v := viper.New()
	v.SetConfigFile(fileName)
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("failed to read config file: %s", err)
		return Config{}
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("failed to unmarshal config file: %s", err)
		return Config{}
	}

	return c
}
