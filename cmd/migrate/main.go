package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please, provide a migration direction: 'up' or 'down'")
	}

	direction := os.Args[1]

	db, err := sql.Open("sqlite3", "./data.db")

	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	src, err := (&file.File{}).Open("cmd/migrate/migrations")

	if err != nil {
		log.Fatal(err.Error())
	}

	mig, err := migrate.NewWithInstance("file", src, "sqlite3", instance)

	if err != nil {
		log.Fatal(err.Error())
	}

	switch direction {

	case "up":
		if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err.Error())
		}

	case "down":
		if err := mig.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err.Error())
		}

	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}
}
