package config

import (
	"log"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Host       string
	Port       string
	DBUser     string
	DBPassword string
	DBname     string
	DBhost     string
	DBSSLMode  bool
}

// Reading configuration from file or environment variables.
func ReadConfig() Config {
	var cfg Config
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	cfg.Host = viper.GetString("SERVER_HOST")
	cfg.Port = viper.GetString("SERVER_PORT")
	cfg.DBUser = viper.GetString("DB_USER")
	cfg.DBPassword = viper.GetString("DB_PASS")
	cfg.DBname = viper.GetString("DB_NAME")
	cfg.DBhost = viper.GetString("DB_HOST")

	log.Println("Environment variables loaded")

	return cfg
}
