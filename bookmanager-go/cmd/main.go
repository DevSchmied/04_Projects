package main

import (
	"bookmanager-go/internal/service"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello bookmanager-go!")

	// Dependency Injection – passing the concrete implementation here
	connector := &service.SQLiteConnector{DBPath: "books.db"}
	db := service.InitDB(connector)

	// Retrieve the underlying *sql.DB to manage the connection lifecycle
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GORM: %v", err)
	}
	defer sqlDB.Close()

	fmt.Printf("Database initialized: %T\n", db)
}
