/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

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
		// Get the sender using the mnemonic in the .env file
		sender, err := blockchain.GetSender(settings.Mnemonic)

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
		if err != nil {
			fmt.Printf("Error creating transaction: %q\n", err)
		}

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
