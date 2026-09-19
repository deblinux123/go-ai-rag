package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func OpenDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "chat.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = on")

	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Database connected")

	return db, nil
}
