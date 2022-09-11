/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
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
		sender, err := blockchain.GetSender(settings.Mnemonic, settings.DerivationPath)
		if err != nil {
			fmt.Printf("Error getting the sender %q\n", err)
			return
		}

		senderAddress, err := blockchain.HexToBech32(sender.PrivKey.PubKey().Address().String())
		if err != nil {
			fmt.Printf("Error converting the sender address %q\n", err)
			return
		}

		granters, err := database.GetGrantersFromDb()
		if err != nil {
			fmt.Println(err)
			return
		}

		grantersToMessage := []messages.ValueToClaim{}
		for _, v := range granters {
			if v.IsValidator == true {
				// Commissions
				// TODO: check if it's needed to send the GetCommission message or if it's autoclaimend when restaking
				res, err := blockchain.GetCommission(v.Validator)
				if err != nil {
					fmt.Printf("Error getting the commission for validator %s, %q\n", v.Validator, err)
					continue
				}

				total := sdkmath.NewIntFromUint64(0)
				for _, t := range res.Commission.Commission {
					if t.Denom == settings.FeeDenom {
						amountToParse := t.Amount
						amount := strings.Split(t.Amount, ".")
						if len(amount) == 2 {
							amountToParse = amount[0]
						}
						val, ok := sdkmath.NewIntFromString(amountToParse)
						if ok != true {
							fmt.Printf("Error parsing commission from %s, %q\n", v.Validator, err)
							continue
						}
						total = total.Add(val)
					}
				}

				if total.LT(sdkmath.NewIntFromUint64(settings.MinReward)) {
					fmt.Printf("NOT enough commission (%d%s) to claim from %s\n", total, settings.FeeDenom, v.Validator)
					continue
				}
				grantersToMessage = append(grantersToMessage, messages.ValueToClaim{
					Granter:     v.Address,
					Validator:   v.Validator,
					Denom:       settings.FeeDenom,
					Amount:      total,
					IsValidator: v.IsValidator,
				})
				fmt.Printf("Claiming commission for %s\n", v.Address)
			} else {
				// Delegations
				res, err := blockchain.GetDistributionRewards(v.Address)
				if err != nil {
					fmt.Printf("Error getting rewards from %s, %q", v.Address, err)
					continue
				}

				total := sdkmath.NewIntFromUint64(0)
				for _, r := range res.Rewards {
					if r.ValidatorAddress == v.Validator {
						for _, t := range r.Reward {
							if t.Denom == settings.FeeDenom {
								amountToParse := t.Amount
								amount := strings.Split(t.Amount, ".")
								if len(amount) == 2 {
									amountToParse = amount[0]
								}
								val, ok := sdkmath.NewIntFromString(amountToParse)
								if ok != true {
									fmt.Printf("Error parsing rewards from %s, %q\n", v.Address, err)
									continue
								}
								total = total.Add(val)
							}
						}
					}
				}

				if total.LT(sdkmath.NewIntFromUint64(settings.MinReward)) {
					fmt.Printf("NOT enough rewards (%s%s) to claim and restake from %s\n", total.String(), settings.FeeDenom, v.Address)
					continue
				}

				grantersToMessage = append(grantersToMessage, messages.ValueToClaim{
					Granter:     v.Address,
					Validator:   v.Validator,
					Denom:       settings.FeeDenom,
					Amount:      total,
					IsValidator: v.IsValidator,
				})

				fmt.Printf("Claiming and restaking %s%s from %s\n", total.String(), settings.FeeDenom, v.Address)
			}
		}

		if len(grantersToMessage) == 0 {
			fmt.Println("Not sending transactions, the granters array is empty!")
			return
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

		fmt.Println(tx)

		// Broadcast the transaction
		// txHash, err := blockchain.Broadcast(tx)
		// if err != nil {
		// 	fmt.Printf("Error broadcasting... %q\n", err)
		// } else {
		// 	fmt.Printf("Transaction included in a block with hash %q\n", txHash)
		// }
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
