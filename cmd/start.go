/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/hanchon-live/autostake-bot/internal/blockchain"
	"github.com/hanchon-live/autostake-bot/internal/util"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethermintcodec "github.com/evmos/ethermint/crypto/codec"
)

var settings util.Config

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "It sends all the restake transactions",
	Long: `It will read the database to get all the granters
It will query the balance for each granter and if it's greater than 0.1 Evmos,
it will claim and restake the total amount`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the private key using the mnemonic in the .env file
		// NOTE: it will use the first wallet, the path is hardcoded
		priv, err := blockchain.CreatePrivateKeyFromMnemonic(settings.Mnemonic)
		if err != nil {
			fmt.Printf("Error creating priv: %q\n", err)
			return
		}

		address, err := blockchain.HexToBech32(priv.PubKey().Address().String())
		if err != nil {
			fmt.Printf("Error getting the address: %q\n", err)
			return
		}

		fmt.Printf("Using the address: %s\n", address)

		// Query the account info to the blockchain
		response, err := blockchain.GetAccountFromBlockchain(address)
		if err != nil {
			fmt.Printf("Error getting account data from the blockchain: %q\n", err)
			return
		}
		sequence, err := strconv.ParseInt(response.Account.BaseAccount.Sequence, 10, 64)
		accountNumber, err2 := strconv.ParseInt(response.Account.BaseAccount.AccountNumber, 10, 64)
		if err != nil || err2 != nil {
			fmt.Println("Error parsing the sequence or account number")
			return
		}

		// Create the sender
		sender := blockchain.Sender{
			Sequence:      uint64(sequence),
			AccountNumber: uint64(accountNumber),
			PrivKey:       priv,
		}

		// Create the proto message
		from, err := blockchain.Bech32StringToAddress("evmos10gu0eudskw7nc0ef48ce9x22sx3tft0s463el3")
		if err != nil {
			fmt.Printf("Error creating from address: %q\n", err)
			return
		}
		to, err := blockchain.Bech32StringToAddress("evmos1urc5gn9x4kvl3sxu4qd9vckfdmtet7shdskm55")
		if err != nil {
			fmt.Printf("Error creating to address: %q\n", err)
		}

		msgSend := bank.NewMsgSend(from, to, blockchain.Uint64ToCoins(int64(42069), settings.FeeDenom))

		// Create the Enconder
		reg := codectypes.NewInterfaceRegistry()
		bank.RegisterInterfaces(reg)
		ethermintcodec.RegisterInterfaces(reg)
		enconder := codec.NewProtoCodec(reg)

		// Create the message
		message := blockchain.Message{
			Msg:      msgSend,
			Enconder: *enconder,
			Fee:      blockchain.Uint64ToCoins(settings.Fee, settings.FeeDenom),
			GasLimit: settings.GasLimit,
			Memo:     settings.Memo,
			ChainId:  settings.ChainId,
		}

		tx, err := blockchain.CreateTransaction(sender, message)

		// Broadcast the transaction
		txHash, err := blockchain.Broadcast(tx)
		if err != nil {
			fmt.Printf("Error broadcasting... %q\n", err)
		} else {
			fmt.Printf("Transaction included in a block with hash %q\n", txHash)
		}

	},
}

func init() {
	config, err := util.LoadConfig()
	if err != nil {
		fmt.Println("Error reading the config, using localnet values!")
	}
	settings = config

	rootCmd.AddCommand(startCmd)
}
