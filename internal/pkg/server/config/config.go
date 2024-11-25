package config

/*
*  This file has been place in its own package rather in the server package
*  to avoid circular import.
*/

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	CertKeyPath string `mapstructure:"cert_key"`
	CertPemPath string `mapstructure:"cert_pub"`

	HttpAddr string `mapstructure:"http_addr"`
	HttpPort string `mapstructure:"http_port"`
	TcpAddr  string `mapstructure:"server_addr"`
	TcpPort  string `mapstructure:"server_port"`

	GlobalDomainName string `mapstructure:"global_domain_name"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"database_path"`
}

// use keycloak
type AuthConfig struct {
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

var GlobalConfig Config

func InitConfig() {
	GlobalConfig = Config{
		App: AppConfig{
			CertKeyPath: "./certs/server.key",
			CertPemPath: "./certs/server.pem",
			HttpAddr:    "0.0.0.0",
			HttpPort:    "8080",
			TcpAddr:     "0.0.0.0",
			TcpPort:     "4040",
		},
		Database: DatabaseConfig{
			Path: "./moon.db",
		},
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/moon/")
	viper.AddConfigPath("$HOME/.config/moon/")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	err := viper.Unmarshal(&GlobalConfig)
	if err != nil {
		log.Fatalf("Unable to unmarshal config file, %v", err)
	}

	if GlobalConfig.App.GlobalDomainName == "" {
		log.Fatalf("'global_domain_name' can't be empty.")
	}
}
