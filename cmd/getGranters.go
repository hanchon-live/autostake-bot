/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	//"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/hanchon-live/autostake-bot/internal/requester"
	"github.com/spf13/cobra"
)

// getGrantersCmd represents the getGranters command
var getGrantersCmd = &cobra.Command{
	Use:   "getGranters",
	Short: "Queries the blockchain to get all the granters and store them in the db.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getGranters called")
		url := "/cosmos/authz/v1beta1/grants/grantee/" + "evmos1u989x5x4vqrkryj8v549yvxpf3yfg9nx3racqn"
		if resp, err := requester.MakeGetRequest("rest", url); err != nil {
			fmt.Println("Failed to get grants")
			return
		} else {
			//	var m authz.QueryGranteeGrantsResponse
			fmt.Println(resp)
		}
	},
}

func init() {
	rootCmd.AddCommand(getGrantersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getGrantersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getGrantersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
