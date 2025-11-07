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
		{Title: "Moby-Dick", Author: "Herman Melville", Year: 1851, Genre: "Adventure", ISBN: "9781503280786", Rating: 4.1, Read: false},
		{Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869, Genre: "Historical", ISBN: "9780199232765", Rating: 4.4, Read: false},
		{Title: "Crime and Punishment", Author: "Fyodor Dostoevsky", Year: 1866, Genre: "Psychological", ISBN: "9780140449136", Rating: 4.7, Read: true},
		{Title: "Brave New World", Author: "Aldous Huxley", Year: 1932, Genre: "Dystopian", ISBN: "9780060850524", Rating: 4.5, Read: true},
		{Title: "The Catcher in the Rye", Author: "J.D. Salinger", Year: 1951, Genre: "Fiction", ISBN: "9780316769488", Rating: 4.0, Read: false},
		{Title: "The Lord of the Rings", Author: "J.R.R. Tolkien", Year: 1954, Genre: "Fantasy", ISBN: "9780544003415", Rating: 4.9, Read: true},
		{Title: "Animal Farm", Author: "George Orwell", Year: 1945, Genre: "Political Satire", ISBN: "9780451526342", Rating: 4.6, Read: true},
		{Title: "Jane Eyre", Author: "Charlotte Brontë", Year: 1847, Genre: "Romance", ISBN: "9780142437209", Rating: 4.5, Read: false},
		{Title: "Wuthering Heights", Author: "Emily Brontë", Year: 1847, Genre: "Tragedy", ISBN: "9780141439556", Rating: 4.3, Read: false},
		{Title: "The Alchemist", Author: "Paulo Coelho", Year: 1988, Genre: "Philosophical", ISBN: "9780061122415", Rating: 4.2, Read: true},
		{Title: "The Da Vinci Code", Author: "Dan Brown", Year: 2003, Genre: "Thriller", ISBN: "9780307474278", Rating: 4.0, Read: false},
		{Title: "Harry Potter and the Sorcerer’s Stone", Author: "J.K. Rowling", Year: 1997, Genre: "Fantasy", ISBN: "9780590353427", Rating: 4.9, Read: true},
		{Title: "The Kite Runner", Author: "Khaled Hosseini", Year: 2003, Genre: "Drama", ISBN: "9781594631931", Rating: 4.8, Read: true},
		{Title: "A Game of Thrones", Author: "George R.R. Martin", Year: 1996, Genre: "Fantasy", ISBN: "9780553593716", Rating: 4.7, Read: false},
		{Title: "The Hunger Games", Author: "Suzanne Collins", Year: 2008, Genre: "Dystopian", ISBN: "9780439023528", Rating: 4.5, Read: true},
		{Title: "The Fault in Our Stars", Author: "John Green", Year: 2012, Genre: "Romance", ISBN: "9780525478812", Rating: 4.4, Read: false},
		{Title: "The Girl with the Dragon Tattoo", Author: "Stieg Larsson", Year: 2005, Genre: "Crime", ISBN: "9780307454546", Rating: 4.6, Read: false},
		{Title: "The Shining", Author: "Stephen King", Year: 1977, Genre: "Horror", ISBN: "9780307743657", Rating: 4.7, Read: false},
		{Title: "Sapiens: A Brief History of Humankind", Author: "Yuval Noah Harari", Year: 2011, Genre: "Nonfiction", ISBN: "9780062316097", Rating: 4.8, Read: true},
		{Title: "The Pragmatic Programmer", Author: "Andrew Hunt & David Thomas", Year: 1999, Genre: "Programming", ISBN: "9780201616224", Rating: 4.9, Read: true},
	}
	for _, b := range books {
		if err := db.Create(&b).Error; err != nil {
			fmt.Println("Error while adding test data:", err)
		}
	}
	fmt.Printf("%d test books added successfully.\n", len(books))
}
