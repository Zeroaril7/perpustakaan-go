package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type envConfig struct {
	AppName           string
	AppPort           string
	BasicAuthUsername string
	BasicAuthPassword string
	MySQLHost         string
	MySQLUsername     string
	MySQLPassword     string
	MySQLDBName       string
	PrivateKey        string
	PublicKey         string
}

var envCfg envConfig

func init() {
	LoadConfig()
}

func LoadConfig() {
	err := godotenv.Load()

	if err != nil {
		println(err.Error())
	}

	envCfg = envConfig{
		AppName:           os.Getenv("APP_NAME"),
		AppPort:           os.Getenv("APP_PORT"),
		BasicAuthUsername: os.Getenv("BASIC_AUTH_USERNAME"),
		BasicAuthPassword: os.Getenv("BASIC_AUTH_PASSWORD"),
		MySQLHost:         os.Getenv("MYSQL_HOST"),
		MySQLUsername:     os.Getenv("MYSQL_USERNAME"),
		MySQLPassword:     os.Getenv("MYSQL_PASSWORD"),
		MySQLDBName:       os.Getenv("MYSQL_DB_NAME"),
		PrivateKey:        os.Getenv("PRIVATE_KEY"),
		PublicKey:         os.Getenv("PUBLIC_KEY"),
	}
}

func (e envConfig) MySQLDSN() (string, string) {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", envCfg.MySQLUsername, envCfg.MySQLPassword, envCfg.MySQLHost, envCfg.MySQLDBName), envCfg.MySQLDBName
}

func Config() *envConfig {
	return &envCfg
}
