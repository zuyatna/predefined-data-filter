package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	httpDelivery "predefined-data-filter/internal/delivery/http"
	"predefined-data-filter/internal/repository/postgres"
	"predefined-data-filter/internal/usecase"
)

func main() {
	// 1. Setup Database Connection
	// These usually come from environment variables.
	connStr := "host=localhost port=5432 user=user password=password dbname=ecommercedb sslmode=disable"

	var db *sql.DB
	var err error

	// Retry connection for docker-compose startup
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Println("Waiting for database...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)
	log.Println("Connected to database successfully")

	// 2. Initialize Layers (Clean Architecture)
	productRepo := postgres.NewProductRepository(db)
	productUseCase := usecase.NewProductUseCase(productRepo)

	// 3. Setup HTTP Router
	mux := http.NewServeMux()
	httpDelivery.NewProductHandler(mux, productUseCase)

	// 4. Start HTTP Server
	serverPort := ":8080"
	log.Printf("Starting server on http://localhost%s\n", serverPort)
	if err := http.ListenAndServe(serverPort, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
