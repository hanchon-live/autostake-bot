/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes local database",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Remove("./autostake-bot.db")
		if err != nil {
			fmt.Println("Couldn't delete the database")
		} else {
			fmt.Println("Database deleted")
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
