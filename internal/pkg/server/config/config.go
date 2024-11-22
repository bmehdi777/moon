package config

import (
	"log"

	"github.com/spf13/viper"
)


type Config struct {
	CertKeyPath string `mapstructure:"cert_key"`
	CertPemPath string `mapstructure:"cert_pub"`

	HttpAddr string `mapstructure:"http_addr"`
	HttpPort string `mapstructure:"http_port"`
	TcpAddr string `mapstructure:"server_addr"`
	TcpPort string `mapstructure:"server_port"`

	DatabasePath string `mapstructure:"database_path"`
}

var GlobalConfig Config

func InitConfig() {
	GlobalConfig = Config{
		CertKeyPath: "./certs/server.key",
		CertPemPath: "./certs/server.pem",
		HttpAddr: "0.0.0.0",
		HttpPort: "8080",
		TcpAddr: "0.0.0.0",
		TcpPort: "4040",
		DatabasePath: "./moon.db",
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
}
