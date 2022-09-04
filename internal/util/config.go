package util

import "github.com/spf13/viper"

type Config struct {
	Rest     []string `mapstructure:"rest"`
	Jrpc     []string `mapstructure:"jrpc"`
	Web3     []string `mapstructure:"web3"`
	Mnemonic string   `mapstructure:"mnemonic"`
}

var defaultConfig = Config{
	Rest:     []string{"http://127.0.0.1:1317"},
	Jrpc:     []string{"http://127.0.0.1:26657"},
	Web3:     []string{"http://127.0.0.1:8545"},
	Mnemonic: "flash local taste power maple fragile pool name file position drop swarm",
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		config = defaultConfig
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		config = defaultConfig
	}

	return
}
