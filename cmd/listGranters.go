/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/hanchon-live/autostake-bot/internal/database"
	"github.com/spf13/cobra"
)

// listGrantersCmd represents the listGranters command
var listGrantersCmd = &cobra.Command{
	Use:   "listGranters",
	Short: "List all the granters stored in the database",
	Run: func(cmd *cobra.Command, args []string) {
		granters, err := database.GetGrantersFromDb()
		if err != nil {
			fmt.Println(err)
		}

		for k, v := range granters {
			fmt.Printf("- %d: Granter->%s. Validator->%s. IsValidator->%t \n", k, v.Address, v.Validator, v.IsValidator)
		}
	},
}

func init() {
	rootCmd.AddCommand(listGrantersCmd)
}
