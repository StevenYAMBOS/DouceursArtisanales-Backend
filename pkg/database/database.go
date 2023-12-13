package database

import (
	"database/sql"
	"fmt"
	"log"

	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	// "github.com/joho/godotenv"
)


var DB *sql.DB

func InitDB() {
    // Pour utiliser dotenv on doit d'abord appeler la méthode godotenv.Load() puis utiliser Getenv après
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Impossible d'accéder au fichier `.env`.")
    }
    dsn := os.Getenv("DATABASE_CREDENTIALS")
    // dsn := "user=stevenyambos password=Zk9c3FzsP94d8K dbname=hmmmdb sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
        fmt.Print("Houla...")
    }
    DB = db 
    fmt.Println("Connexion à la base de donnée étabit !")
}