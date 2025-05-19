package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		URL string
	}
	Rabbit struct {
		User     string
		Password string
		Host     string
	}
	Weather struct {
		GRPC struct {
			Host string
			Port string
		}
	}
	HTTP struct {
		Port    string
		BaseURL string `mapstructure:"base_url"`
	}
	Notification struct {
		Interval string
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
		fmt.Printf("Config file not found: %v\n", err)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Sprintf("Unable to decode into struct: %v", err))
	}
}
