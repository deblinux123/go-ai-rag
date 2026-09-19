package main

import "fmt"

func main() {
	db, err := OpenDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	repo := NewChatRepository(db)

	err = repo.CreateTable()

	if err != nil {
		panic(err)
	}

	fmt.Println("Table is ready.")

	id, err := repo.CreateChat("Learning Golang + Sqlite")

	if err != nil {
		panic(err)
	}

	fmt.Println("Created Chat with ID:", id)

	id, err = repo.CreateChat("Learning AI Desktop application")

	if err != nil {
		panic(err)
	}

	fmt.Println("Created Chat with ID:", id)

	id, err = repo.CreateChat("This is the best way to leatrn new Thing")

	if err != nil {
		panic(err)
	}

	fmt.Println("Created Chat with ID:", id)

	fmt.Print("\n--- All Chats ---\n")

	chats, err := repo.GetChats()

	if err != nil {
		panic(err)
	}

	for _, chat := range chats {
		fmt.Printf("ID: %d | Title: %s\n", chat.ID, chat.Title)
	}

	fmt.Print("\n--- Get one chat ---\n")

	chat, err := repo.GetChat(1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Found Chat: ")
	fmt.Printf("ID: %d | Title: %s\n", chat.ID, chat.Title)

	fmt.Print("\n--- Update ---\n")

	err = repo.UpdateChat(1, "LEARNING NEW GOLANG THING")

	if err != nil {
		panic(err)
	}

	chat, err = repo.GetChat(1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Found Chat after update: ")
	fmt.Printf("ID: %d | Title: %s\n", chat.ID, chat.Title)

	fmt.Print("\n--- Delete ---\n")
	err = repo.DeleteChat(3)

	if err != nil {
		panic(err)
	}

	fmt.Println("Chat Deleted.")

	fmt.Print("\n--- Final Chats ---\n")

	chats, err = repo.GetChats()

	if err != nil {
		panic(err)
	}

	for _, chat := range chats {
		fmt.Printf(
			"ID: %d | Title: %s\n",
			chat.ID,
			chat.Title,
		)
	}

}
