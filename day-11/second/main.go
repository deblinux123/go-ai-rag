package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Chat struct {
	Id    int
	Title string
}

func main() {
	db, err := sql.Open("sqlite", "chat.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println("Database Opend.")

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connection successfuly.")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL
		) 
	`)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table Chats Created.")

	result, err := db.Exec(
		"INSERT INTO chats (title) VALUES (?)",
		"Learning golang",
	)

	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Create chat with id:", id)

	result, err = db.Exec(
		"INSERT INTO chats (title) VALUES (?)",
		"learging sqlite",
	)

	if err != nil {
		log.Fatal(err)
	}

	id, err = result.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created chat with id:", id)

	result, err = db.Exec(
		"INSERT INTO chats (title) VALUES (?)",
		"Building AI Chat",
	)

	if err != nil {
		log.Fatal(err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\n----- All Chats -----\n")

	rows, err := db.Query(
		"SELECT id, title FROM chats ORDER BY id",
	)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var chat Chat

		err := rows.Scan(
			&chat.Id, &chat.Title,
		)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"ID: %d | Title: %s\n",
			chat.Id,
			chat.Title,
		)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Print("\n----- Find Chat -----\n")

	var chat Chat

	err = db.QueryRow(
		"SELECT id, title FROM chats WHERE id = ?",
		2,
	).Scan(&chat.Id, &chat.Title)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found Chat: ID=%d | Title=%s\n", chat.Id, chat.Title)

	fmt.Print("\n--- Update Chat ---\n")

	result, err = db.Exec(
		"UPDATE chats SET title = ? WHERE id = ?",
		"Learning go + Sqlite",
		1,
	)

	if err != nil {
		log.Fatal(err)
	}

	affected, err := result.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Rows Updated: ", affected)

	err = db.QueryRow(
		"SELECT id, title FROM chats WHERE id = ?",
		1,
	).Scan(&chat.Id, &chat.Title)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"Updated chat: ID=%d | Title=%s\n",
		chat.Id,
		chat.Title,
	)

	fmt.Print("\n--- Delete Chat ---\n")

	result, err = db.Exec(
		"DELETE FROM chats WHERE ID = ?",
		2,
	)

	if err != nil {
		log.Fatal(err)
	}

	affected, err = result.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Rows deleted:", affected)

	fmt.Print("\n--- Chats After Delete ---\n")

	rows, err = db.Query("SELECT id, title FROM chats ORDER BY ID")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var chat Chat

		err := rows.Scan(
			&chat.Id,
			&chat.Title,
		)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"ID: %d | Title: %s\n",
			chat.Id,
			chat.Title,
		)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}
