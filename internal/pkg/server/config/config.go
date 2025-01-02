package config

/*
*  This file has been place in its own package rather in the server package
*  to avoid circular import.
 */

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

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
	Realm string `mapstructure:"realm"`
	BaseURL string `mapstructure:"base_url"`
	Algorithm string `mapstructure:"algorithm"`
	Audience string `mapstructure:"audience"`
}
type RealmConfig struct {
	Realm string `json:"realm"`
	PublicKey string `json:"public_key"`
	TokenService string `json:"token-service"`
	AccountService string `json:"account-service"`
	TokenNotBefore int `json:"tokens-not-before"`
}

type Config struct {
	App      AppConfig `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig `mapstructure:"auth"`

	RealmConfig RealmConfig
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

	verifyApp()
	verifyAuth()

	GlobalConfig.RealmConfig = getPublicKey()
}


func verifyApp() {
	if GlobalConfig.App.GlobalDomainName == "" {
		log.Fatalf("'app.global_domain_name' can't be empty.")
	}
}

func verifyAuth() {
	if GlobalConfig.Auth.Realm == "" {
		log.Fatalf("'auth.realm' can't be empty.")
	}
	if GlobalConfig.Auth.BaseURL == "" {
		log.Fatalf("'auth.base_url' can't be empty.")
	}
	if GlobalConfig.Auth.Algorithm == "" {
		log.Fatalf("'auth.algorithm' can't be empty.")
	}
	if GlobalConfig.Auth.Audience == "" {
		log.Fatalf("'auth.audience' can't be empty.")
	}
}

func getPublicKey() RealmConfig {
	response, err := http.Get(GlobalConfig.Auth.BaseURL + "/realms/" + GlobalConfig.Auth.Realm)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if response.StatusCode != 200 {
		log.Fatalf("Error : Keycloak %v", response.Status)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var realmConfig RealmConfig
	err = json.Unmarshal(body, &realmConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return realmConfig
}
