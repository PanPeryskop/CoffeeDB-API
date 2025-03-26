package db

import (
    "database/sql"
    "fmt"
    "os"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
    // Wczytanie zmiennych środowiskowych z pliku .env
    if err := godotenv.Load(); err != nil {
        fmt.Println("Nie można wczytać pliku .env, korzystam ze zmiennych systemowych")
    }

    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbName)

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }

    if err := DB.Ping(); err != nil {
        return err
    }

    fmt.Println("Pomyślnie połączono z bazą danych")
    return nil
}