package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "chat.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connection successful!")

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS chats (
			id 	INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table created.")

	result, err := db.Exec(
		"INSERT INTO chats (title) VALUES (?)",
		"this is my first chat that save to database",
	)

	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Chat created with id: ", id)

	rows, err := db.Query("SELECT id, title FROM chats 	WHERE id > 1")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var title string

		err := rows.Scan(&id, &title)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d - %s\n", id, title)
	}

	var newId int
	var title string

	err = db.QueryRow(
		"SELECT id, title FROM chats WHERE id=?",
		3,
	).Scan(&newId, &title)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
	fmt.Println("Found Chat:")
	fmt.Printf("%d : %s\n", newId, title)

	newResult, err := db.Exec(
		"UPDATE chats SET title = ? WHERE id = ?",
		"Learning SQLITE",
		4,
	)

	if err != nil {
		log.Fatal(err)
	}

	affected, err := newResult.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Rows Updated:", affected)

	err = db.QueryRow(
		"SELECT id, title FROM chats WHERE id = ?",
		4,
	).Scan(&newId, &title)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Updated Chat:", newId, title)
}
