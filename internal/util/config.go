package util

import "github.com/spf13/viper"

type Config struct {
	Rest     []string `mapstructure:"rest"`
	Jrpc     []string `mapstructure:"jrpc"`
	Web3     []string `mapstructure:"web3"`
	Mnemonic string   `mapstructure:"mnemonic"`
	Fee      int64    `mapstructure:"fee"`
	FeeDenom string   `mapstructure:"feedenom"`
	GasLimit uint64   `mapstructure:"gaslimit"`
	ChainId  string   `mapstructure:"chainid"`
	Memo     string   `mapstructure:"memo"`
}

var defaultConfig = Config{
	Rest:     []string{"http://127.0.0.1:1317"},
	Jrpc:     []string{"http://127.0.0.1:26657"},
	Web3:     []string{"http://127.0.0.1:8545"},
	Mnemonic: "flash local taste power maple fragile pool name file position drop swarm",
	Fee:      int64(100000000),
	FeeDenom: "aevmos",
	GasLimit: uint64(150000),
	ChainId:  "evmos_9000-1",
	Memo:     "Hanchon restake",
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
