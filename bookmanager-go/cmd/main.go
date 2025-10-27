package main

import (
	"bookmanager-go/internal/model"
	"bookmanager-go/internal/server"
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

	// Run AutoMigrate to create the Book table if it doesn't exist
	if err := db.AutoMigrate(&model.Book{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	serverAddress := "localhost:8080"                                                       // Server listening address
	staticRoute := "/static"                                                                // URL route for static files
	staticPath := "./internal/view/static"                                                  // Local folder for static files
	templatePath := "internal/view/templates/**/*.html"                                     // HTML templates location
	appServer := server.NewServer(db, serverAddress, templatePath, staticRoute, staticPath) // Initialize server with dependencies
	if err := appServer.Start(); err != nil {                                               // Start web server and handle startup errors
		log.Fatalf("Server failed to start: %v", err)
	}
}
