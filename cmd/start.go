/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/hanchon-live/autostake-bot/internal/blockchain"
	"github.com/hanchon-live/autostake-bot/internal/database"
	"github.com/hanchon-live/autostake-bot/internal/messages"
	"github.com/hanchon-live/autostake-bot/internal/util"
	"github.com/spf13/cobra"
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
		sender, _ := blockchain.GetSender(settings.Mnemonic)

		senderAddress, _ := blockchain.HexToBech32(sender.PrivKey.PubKey().Address().String())

		granters, err := database.GetGrantersFromDb()
		if err != nil {
			fmt.Println(err)
			return
		}

		grantersToMessage := []messages.ValueToClaim{}
		for _, v := range granters {
			grantersToMessage = append(grantersToMessage, messages.ValueToClaim{
				Granter:   v.Address,
				Validator: v.Validator,
				Denom:     settings.FeeDenom,
				// TODO: read amount from blockchain
				Amount: 1,
			})
		}

		// Create the proto message
		msg, encoder, err := messages.CreateMessageExec(senderAddress, grantersToMessage)
		if err != nil {
			fmt.Printf("Error creating message send: %q\n", err)
			return
		}

		// Enconde new message
		message := messages.NewMessage(
			&msg,
			encoder,
			settings.Fee,
			settings.FeeDenom,
			settings.GasLimit,
			settings.Memo,
			settings.ChainId,
		)

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
