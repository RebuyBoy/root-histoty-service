package config

import "os"

type DBConfig struct {
	Driver   string
	User     string
	Password string
	Database string
	Port     string
}

type ServerConfig struct {
	Port       string
	SecretWord string
}

func GetDBConfig() DBConfig {
	dbConfig := DBConfig{
		Driver:   os.Getenv("DB_DRIVER"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USERNAME"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	}
	return dbConfig
}

func GetServerConfig() ServerConfig {
	serverConfig := ServerConfig{
		Port:       os.Getenv("SERVER_PORT"),
		SecretWord: os.Getenv("SECRET_WORD"),
	}
	return serverConfig
}
