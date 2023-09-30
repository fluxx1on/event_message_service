package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"gitlab.com/fluxx1on_group/event_message_service/internal/config"
)

func main() {
	migrate := flag.Bool("m", false, "Run database migration")
	flag.Parse()

	if *migrate {
		Migrate()
	} else {
		Drop()
	}
}

func Migrate() {
	filename := "migrations/db.sql"

	conn, err := pgx.Connect(context.Background(), config.NewDB().URL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	dbFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to open %v file %v", filename, err)
	}

	migrationSQL, err := io.ReadAll(dbFile)
	if err != nil {
		log.Fatalf("Unable to read %v file: %v", filename, err)
	}

	_, err = conn.Exec(context.Background(), string(migrationSQL))
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migrations successfully applied.")
}

func Drop() {
	filename := "migrations/drop_db.sql"

	conn, err := pgx.Connect(context.Background(), config.NewDB().URL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	dbFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to open %v file %v", filename, err)
	}

	migrationSQL, err := io.ReadAll(dbFile)
	if err != nil {
		log.Fatalf("Unable to read %v file: %v", filename, err)
	}

	_, err = conn.Exec(context.Background(), string(migrationSQL))
	if err != nil {
		log.Fatalf("Tables dropping failed: %v", err)
	}

	fmt.Println("Tables successfully dropped.")
}
