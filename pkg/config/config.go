package config

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DBPort     string `mapstructure:"DBPort"`
	DBHost     string `mapstructure:"DBHOST"`
	DBName     string `mapstructure:"DBNAME"`
	DBPassword string `mapstructure:"DBPASSWORD"`
}

var envs = []string{
	"DBNAME", "DBPASSWORD",
}

func LoadConfig() (Config, error) {
	var cfg Config
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return cfg, err
		}
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	if err := validator.New().Struct(&cfg); err != nil {
		return cfg, err
	}

	LoadEnv()

	return cfg, nil
}
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading Env File")
	}

}
