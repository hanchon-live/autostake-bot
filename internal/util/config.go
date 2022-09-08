package util

import "github.com/spf13/viper"

type Config struct {
	Rest               []string `mapstructure:"rest"`
	Mnemonic           string   `mapstructure:"mnemonic"`
	DerivationPath     string   `mapstructure:"derivationpath"`
	Fee                int64    `mapstructure:"fee"`
	FeeDenom           string   `mapstructure:"feedenom"`
	GasLimit           uint64   `mapstructure:"gaslimit"`
	ChainId            string   `mapstructure:"chainid"`
	Memo               string   `mapstructure:"memo"`
	Validator          string   `mapstructure:"validator"`
	MultipleValidators bool     `mapstructure:"multiplevalidators"`
	GranteeWallet      string   `mapstructure:"granteewallet"`
	MinReward          uint64   `mapstructure:"minreward"`
}

var defaultConfig = Config{
	Rest:               []string{"http://127.0.0.1:1317"},
	Mnemonic:           "flash local taste power maple fragile pool name file position drop swarm",
	DerivationPath:     "m/44'/60'/0'/0/0",
	Fee:                int64(100000000),
	FeeDenom:           "aevmos",
	GasLimit:           uint64(150000),
	ChainId:            "evmos_9000-1",
	Memo:               "Hanchon restake",
	Validator:          "evmosvaloper1nm5uh2q85h9vylzs6uuvje4cscz4dcew8cawss",
	MultipleValidators: false,
	GranteeWallet:      "evmos10gu0eudskw7nc0ef48ce9x22sx3tft0s463el3",
	MinReward:          uint64(100000000),
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
