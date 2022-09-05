/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hanchon-live/autostake-bot/internal/requester"
	"github.com/hanchon-live/autostake-bot/types/responses"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// updateGrantersCmd represents the getGranters command
var updateGrantersCmd = &cobra.Command{
	Use:   "updateGranters",
	Short: "Queries the blockchain to get all the granters and store them in the db.",
	Run: func(cmd *cobra.Command, args []string) {
		// Open the db
		db, err := sql.Open("sqlite3", "./autostake-bot.db")
		if err != nil {
			fmt.Printf("Error opening database: %q", err)
			return
		}

		defer db.Close()

		// Get all the wallets
		url := "/cosmos/authz/v1beta1/grants/grantee/" + settings.GranteeWallet
		if resp, err := requester.MakeGetRequest("rest", url); err != nil {
			fmt.Println("Failed to get grants")
			return
		} else {
			m := &responses.GrantsReponse{}

			err = json.Unmarshal([]byte(resp), m)
			if err != nil {
				fmt.Printf("Error decoding response: %q", err)
			}

			tx, err := db.Begin()
			if err != nil {
				fmt.Printf("Error creating transaction: %q", err)
			}

			stmt, err := tx.Prepare("insert into delegators(name, address, isvalidator, validator, maxamount) values(?, ?,?,?,?)")
			if err != nil {
				fmt.Printf("Error preparing transaction: %q", err)
			}
			defer stmt.Close()

			for _, grant := range m.Grants {
				// Store it to the database
				if grant.Authorization.Value == "/cosmos.staking.v1beta1.MsgDelegate" {
					_, err = stmt.Exec("", grant.Granter, false, settings.Validator, 0)
					if err != nil && err.Error() == "UNIQUE constraint failed: delegators.address" {
						fmt.Println("=", grant.Grantee, "already stored in db.")
					} else if err != nil {
						fmt.Printf("Error executing transaction: %q", err)
					} else {
						fmt.Println("+", grant.Grantee, "stored in db.")
					}
				}
			}

			err = tx.Commit()
			if err != nil {
				fmt.Printf("Error commiting transaction: %q", err)
			}
		}

		fmt.Println("Granters added to database")
	},
}

func init() {
	rootCmd.AddCommand(updateGrantersCmd)
}
