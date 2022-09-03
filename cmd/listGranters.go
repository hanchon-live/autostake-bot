/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
	"fmt"

	"github.com/hanchon-live/autostake-bot/internal/wallet"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// listGrantersCmd represents the listGranters command
var listGrantersCmd = &cobra.Command{
	Use:   "listGranters",
	Short: "List all the granters stored in the database",
	Run: func(cmd *cobra.Command, args []string) {
		// Open the db
		db, err := sql.Open("sqlite3", "./autostake-bot.db")
		if err != nil {
			fmt.Printf("Error opening database: %q", err)
			return
		}

		defer db.Close()

		rows, err := db.Query("select id, address, validator from delegators")

		if err != nil {
			fmt.Printf("Error creating the query to the database: %q", err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var address string
			var validator string
			err = rows.Scan(&id, &address, &validator)
			if err != nil {
				fmt.Printf("Error getting the row information: %q", err)
			}

			fmt.Printf("- %d: Granter->%s. Validator->%s \n", id, address, validator)
		}
		err = rows.Err()
		if err != nil {
			fmt.Printf("Row error: %q", err)
		}

		_, _, _ = wallet.GetWallet()

	},
}

func init() {
	rootCmd.AddCommand(listGrantersCmd)
}
