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
	Weather struct {
		Api struct {
			Key string
		}
	}
	Port string
}

var Conf Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Print("Config file not found...")
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Sprintf("Unable to decode into struct %e", err))
	}
}
