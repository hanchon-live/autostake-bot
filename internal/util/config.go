package util

import "github.com/spf13/viper"

type Config struct {
	Rest []string `mapstructure:"rest"`
	Jrpc []string `mapstructure:"jrpc"`
	Web3 []string `mapstructure:"web3"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
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
