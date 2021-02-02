package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func ConfigServer(cfgFile string) {

	if cfgFile == "" {
		panic(errors.New("empty config file"))
	}

	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

}
