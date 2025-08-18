package main

import (
	"database/sql"
	"log"

	"github.com/shekharsikku/go-gin-rest-api/internal/database"
	"github.com/shekharsikku/go-gin-rest-api/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	path := env.GetEnvString("DB_PATH", "./tmp/data.db")
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	models := database.NewModels(db)

	app := &application{
		port:      env.GetEnvInt("PORT", 4000),
		jwtSecret: env.GetEnvString("JWT_SECRET", "secret@gin1.10"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err.Error())
	}
}
