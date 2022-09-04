package blockchain

import (
	"fmt"
	"strconv"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

type Sender struct {
	Sequence      uint64
	AccountNumber uint64
	PrivKey       ethsecp256k1.PrivKey
}

func GetSender(mnemonic string) (Sender, error) {
	// NOTE: it will use the first wallet, the path is hardcoded
	priv, err := CreatePrivateKeyFromMnemonic(mnemonic)
	if err != nil {
		return Sender{}, fmt.Errorf("Error creating priv: %q", err)
	}

	address, err := HexToBech32(priv.PubKey().Address().String())
	if err != nil {
		return Sender{}, fmt.Errorf("Error getting the address: %q", err)
	}

	fmt.Printf("Using the address: %s\n", address)

	// Query the account info to the blockchain
	response, err := GetAccountFromBlockchain(address)
	if err != nil {
		return Sender{}, fmt.Errorf("Error getting account data from the blockchain: %q", err)
	}
	sequence, err := strconv.ParseInt(response.Account.BaseAccount.Sequence, 10, 64)
	accountNumber, err2 := strconv.ParseInt(response.Account.BaseAccount.AccountNumber, 10, 64)
	if err != nil || err2 != nil {
		return Sender{}, fmt.Errorf("Error parsing the sequence or account number: %q", err)
	}

	// Create the sender
	return Sender{
		Sequence:      uint64(sequence),
		AccountNumber: uint64(accountNumber),
		PrivKey:       priv,
	}, nil
}
