package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"

	"github.com/navyamani-ch/Account-Management/internal/handlers"
	"github.com/navyamani-ch/Account-Management/internal/services"
	"github.com/navyamani-ch/Account-Management/internal/stores"
)

func main() {
	r := mux.NewRouter()

	db := runMigrations()
	defer db.Close()

	accountStore := stores.NewAccountStore(db)
	transactionStore := stores.NewTransactionStore(db)

	accountService := services.NewAccountService(accountStore)
	transactionService := services.NewTransactionService(transactionStore, accountStore)

	accountHandler := handlers.NewAccountHandler(accountService)
	trnsactionHandler := handlers.NewTransactionHandler(transactionService)

	r.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{account_id}", accountHandler.GetAccount).Methods("GET")
	r.HandleFunc("/transactions", trnsactionHandler.CreateTransaction).Methods("POST")

	log.Println("Server started on :8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func runMigrations() *sql.DB {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, sslMode)
	// Connect using database/sql
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	// Initialize the migrate PostgreSQL driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Build the full path to the migrations folder
	migrationsPath := "file://" + filepath.Join(wd, "migrations")

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrate instance: %v", err)
	}

	// Apply migrations
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply.")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		log.Println("Database migrated successfully.")
	}

	return db
}
