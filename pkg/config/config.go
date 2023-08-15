package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DbHost            string
	DbPort            string
	DbUser            string
	DbPass            string
	DbName            string
	ServerPort        string
	DBMigrationFolder string
}

func LoadConfig() (*Config, error) {
	// Load Config
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AutomaticEnv()
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	v.AddConfigPath("./pkg/config")
	v.AddConfigPath("../pkg/config")
	v.AddConfigPath("../../pkg/config")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("[WARN] config file not found, will use environment variables instead.")
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	// Set Defaults
	v.SetDefault("MYSQL_HOST", "127.0.0.1")
	v.SetDefault("MYSQL_PORT", 3306)
	v.SetDefault("MYSQL_USER", "inventory_admin")
	v.SetDefault("MYSQL_PASSWORD", "password")
	v.SetDefault("MYSQL_DATABASE", "inventory")
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("DB_MIGRATION_FOLDER", "file://db/migrations")

	// Load Config Values
	c := Config{
		DbHost:            v.GetString("MYSQL_HOST"),
		DbPort:            v.GetString("MYSQL_PORT"),
		DbUser:            v.GetString("MYSQL_USER"),
		DbPass:            v.GetString("MYSQL_PASSWORD"),
		DbName:            v.GetString("MYSQL_DATABASE"),
		ServerPort:        v.GetString("SERVER_PORT"),
		DBMigrationFolder: v.GetString("DB_MIGRATION_FOLDER"),
	}

	return &c, nil
}
