package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "chat.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	fmt.Println("Database connected!")
}
