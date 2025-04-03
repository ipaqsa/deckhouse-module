package config

import (
	"github.com/spf13/viper"

	"github.com/x/module/internal/service/keyvalue"
	"github.com/x/module/internal/service/roller"
	"github.com/x/module/internal/service/sayer"
)

type Config struct {
	Roller   roller.Config   `mapstructure:"roller"`
	Sayer    sayer.Config    `mapstructure:"sayer"`
	KeyValue keyvalue.Config `mapstructure:"kv"`
}

func Load() Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	conf := new(Config)
	if err := viper.UnmarshalKey("config", conf); err != nil {
		panic(err)
	}

	return *conf
}
