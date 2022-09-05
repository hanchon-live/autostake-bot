package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type GranterFromDb struct {
	Address     string
	Validator   string
	IsValidator bool
}

func GetGrantersFromDb() ([]GranterFromDb, error) {
	// Open the db
	db, err := sql.Open("sqlite3", "./autostake-bot.db")
	if err != nil {
		return []GranterFromDb{}, fmt.Errorf("Error opening database: %q", err)
	}

	defer db.Close()

	rows, err := db.Query("select id, address, validator, isvalidator from delegators")

	if err != nil {
		return []GranterFromDb{}, fmt.Errorf("Error creating the query to the database: %q", err)
	}
	defer rows.Close()

	res := []GranterFromDb{}

	for rows.Next() {
		var id int
		var address string
		var validator string
		var isValidator bool
		err = rows.Scan(&id, &address, &validator, &isValidator)
		if err != nil {
			fmt.Printf("Error getting the row information: %q", err)
		}
		res = append(res, GranterFromDb{Address: address, Validator: validator, IsValidator: isValidator})
	}
	err = rows.Err()

	if err != nil {
		return []GranterFromDb{}, fmt.Errorf("Row error: %q", err)
	}

	return res, nil
}
