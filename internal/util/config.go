package util

import "github.com/spf13/viper"

type Config struct {
	Rest     []string `mapstructure:"rest"`
	Jrpc     []string `mapstructure:"jrpc"`
	Web3     []string `mapstructure:"web3"`
	Mnemonic string   `mapstructure:"mnemonic"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
