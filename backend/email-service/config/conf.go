package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Rabbit struct {
		User     string
		Password string
		Host     string
	}
	Email struct {
		SMTPHost     string
		SMTPUser     string
		SMTPPassword string
		FromAddress  string
	}
	HTTP struct {
		BaseURL string
	}
}

var Conf Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Config file not found, using environment variables")
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Sprintf("Unable to decode config: %v", err))
	}
}
