package wallet

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/hanchon-live/autostake-bot/internal/util"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

var settings util.Config

func init() {
	config, err := util.LoadConfig()
	if err != nil {
		panic("Error reading the config")
	}
	settings = config
}

func GetWallet() (hdwallet.Wallet, accounts.Account, error) {
	wallet, err := hdwallet.NewFromMnemonic(settings.Mnemonic)
	if err != nil {
		return hdwallet.Wallet{}, accounts.Account{}, err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return hdwallet.Wallet{}, accounts.Account{}, err
	}

	fmt.Println(account.Address.String())
	return *wallet, account, nil
}
