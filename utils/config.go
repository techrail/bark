package utils

import (
	"os"
)

type AppConfig struct {
	DBName     string
	DBUsername string
	DBPassword string
	SSLMode    string
	APPPort    string
}

func LoadConfig() AppConfig {
	return AppConfig{
		DBName:     os.Getenv("DBNAME"),
		DBUsername: os.Getenv("DBUSERNAME"),
		DBPassword: os.Getenv("DBPASSWORD"),
		SSLMode:    os.Getenv("SSLMODE"),
		APPPort:    os.Getenv("APPPORT"),
	}
}
