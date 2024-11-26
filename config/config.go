package config

import(
	"context"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func InitConfig(){
	// Load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .en file: %v", err)
	}
	fmt.Println(".env variables loaded successfully.")
}

func InitDB()(*pgx.Conn){
	// Constructing the connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
		os.Getenv("DB_USER"), 
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_HOST"), 
		os.Getenv("DB_PORT"), 
		os.Getenv("DB_NAME"))

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	fmt.Println("Connected to the database.")

	// Set to the global DB variable
	db = conn
	return conn
}

func CloseDB(){
	if db != nil {
		err := db.Close(context.Background())
		if err != nil {
			log.Printf("Error closing the database connection: %v", err)
		} else {
			fmt.Println("Database connection closed.")
		}
	}
}