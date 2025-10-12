package main

import (
	"bookmanager-go/internal/service"
	"fmt"
)

func main() {
	fmt.Println("Hello bookmanager-go!")

	// Dependency Injection – passing the concrete implementation here
	connector := &service.SQLiteConnector{DBPath: "books.db"}
	db := service.InitDB(connector)

	fmt.Printf("Database initialized: %T\n", db)
}
