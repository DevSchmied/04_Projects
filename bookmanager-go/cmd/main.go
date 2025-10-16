package main

import (
	"bookmanager-go/internal/controller"
	"bookmanager-go/internal/model"
	"bookmanager-go/internal/server"
	"bookmanager-go/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
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

	router := gin.Default() // Initialize Gin router
	bookController := controller.BookController{DB: db}
	// Start the web server with injected router, controller, and address
	if err := server.StartWebServer(router, bookController, "localhost:8080"); err != nil {
		log.Panicf("Server failed: %v", err)
	}

}
