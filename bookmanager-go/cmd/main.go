package main

import (
	"bookmanager-go/internal/model"
	"bookmanager-go/internal/server"
	"bookmanager-go/internal/service"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func main() {
	fmt.Println("Hello bookmanager-go!")

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

	// addTestData(db)

	serverAddress := "localhost:8080"                                                       // Server listening address
	staticRoute := "/static"                                                                // URL route for static files
	staticPath := "./internal/view/static"                                                  // Local folder for static files
	templatePath := "internal/view/templates/**/*.html"                                     // HTML templates location
	appServer := server.NewServer(db, serverAddress, templatePath, staticRoute, staticPath) // Initialize server with dependencies
	if err := appServer.Start(); err != nil {                                               // Start web server and handle startup errors
		log.Fatalf("Server failed to start: %v", err)
	}
}

func addTestData(db *gorm.DB) {
	books := []model.Book{
		{Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Year: 1925, Genre: "Novel", ISBN: "9780743273565", Rating: 4.3, Read: true},
		{Title: "To Kill a Mockingbird", Author: "Harper Lee", Year: 1960, Genre: "Fiction", ISBN: "9780061120084", Rating: 4.8, Read: true},
		{Title: "1984", Author: "George Orwell", Year: 1949, Genre: "Dystopian", ISBN: "9780451524935", Rating: 4.6, Read: true},
		{Title: "Pride and Prejudice", Author: "Jane Austen", Year: 1813, Genre: "Romance", ISBN: "9780141439518", Rating: 4.5, Read: false},
		{Title: "The Hobbit", Author: "J.R.R. Tolkien", Year: 1937, Genre: "Fantasy", ISBN: "9780547928227", Rating: 4.7, Read: false},
	}
	for _, b := range books {
		if err := db.Create(&b).Error; err != nil {
			fmt.Println("Error while adding test data:", err)
		}
	}
	fmt.Printf("%d test books added successfully.\n", len(books))
}
