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

type DatabaseDriver string

const (
	DRIVER_SQLITE   DatabaseDriver = "sqlite"
	DRIVER_POSTGRES DatabaseDriver = "postgres"
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
	Driver DatabaseDriver `mapstructure:"driver"`

	// sqlite configuration
	SqlitePath string `mapstructure:"sqlite_path"`

	// postgres configuration
	PostgresUser     string `mapstructure:"postgres_user"`
	PostgresPassword string `mapstructure:"postgres_password"`
	PostgresDbName   string `mapstructure:"postgres_dbname"`
	PostgresPort     string `mapstructure:"postgres_port"`
}

// use keycloak
type AuthConfig struct {
	Realm     string `mapstructure:"realm"`
	BaseURL   string `mapstructure:"base_url"`
	Algorithm string `mapstructure:"algorithm"`
	Audience  string `mapstructure:"audience"`
}
type RealmConfig struct {
	Realm          string `json:"realm"`
	PublicKey      string `json:"public_key"`
	TokenService   string `json:"token-service"`
	AccountService string `json:"account-service"`
	TokenNotBefore int    `json:"tokens-not-before"`
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`

	RealmConfig RealmConfig
}

var GlobalConfig Config

func InitConfig() {
	GlobalConfig = Config{
		App: AppConfig{
			CertKeyPath: "./assets/server/certs/server.key",
			CertPemPath: "./assets/server/certs/server.pem",
			HttpAddr:    "0.0.0.0",
			HttpPort:    "8080",
			TcpAddr:     "0.0.0.0",
			TcpPort:     "4040",
		},
		Database: DatabaseConfig{
			SqlitePath: "./moon.db",
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
	verifyDatabase()
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

func verifyDatabase() {
	switch GlobalConfig.Database.Driver {
	case DRIVER_SQLITE:
		if GlobalConfig.Database.SqlitePath == "" {
			log.Fatalf("'database.sqlite_path' can't be empty if 'database.driver' is set to 'sqlite'.")
		}
		break
	case DRIVER_POSTGRES:
		if GlobalConfig.Database.PostgresUser == "" {
			log.Fatalf("'database.postgres_user' can't be empty if 'database.driver' is set to 'postgres'.")
		}
		if GlobalConfig.Database.PostgresPassword == "" {
			log.Fatalf("'database.postgres_password' can't be empty if 'database.driver' is set to 'postgres'.")
		}
		if GlobalConfig.Database.PostgresDbName == "" {
			log.Fatalf("'database.postgres_dbname' can't be empty if 'database.driver' is set to 'postgres'.")
		}
		if GlobalConfig.Database.PostgresPort == "" {
			log.Fatalf("'database.postgres_port' can't be empty if 'database.driver' is set to 'postgres'.")
		}
		break
	default:
		log.Fatalf("'database.driver' can't be empty.")
		break
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
