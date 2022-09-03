/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates the database and init the tables",
	Run: func(cmd *cobra.Command, args []string) {

		db, err := sql.Open("sqlite3", "./autostake-bot.db")
		if err != nil {
			fmt.Printf("Error creating/opening database: %q", err)
			return
		}
		defer db.Close()

		sqlStmt := `
	create table if not exists delegators (
        id integer not null primary key,
        name text,
        address text,
        isvalidator bool,
        validator text,
        maxamount text
    );
	`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Printf("Error executing the table creation: %q", err)
			return
		}
		fmt.Println("Database initialized")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
