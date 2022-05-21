package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4/database"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func ensureVersionTable(conn *sql.DB) error {
	// check if migration table exists
	var count int
	query := `SELECT COUNT(1) FROM information_schema.tables WHERE table_name = $1 AND table_schema = (SELECT current_schema()) LIMIT 1`
	if err := conn.QueryRowContext(context.Background(), query, "schema_migrations").Scan(&count); err != nil {
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}
	if count == 1 {
		return nil
	}

	// if not, create the empty migration table
	query = `CREATE TABLE "schema_migrations" (version bigint not null primary key, dirty boolean not null)`
	if _, err := conn.ExecContext(context.Background(), query); err != nil {
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}
	return nil
}

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	if !strings.Contains(dbUrl, "test") {
		log.Fatalln("invalid database")
	}

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatalln("could not connect to database", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalln("could not get driver", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./storages/postgres/migrations/",
		"tests",
		driver,
	)
	if err != nil {
		log.Fatalln("could not get migrate instance", err)
	}
	err = m.Drop()
	if err != nil {
		log.Fatalln("could not drop database", err)
	}
	err = ensureVersionTable(db)
	if err != nil {
		log.Fatalln("could not ensure version table")
	}
	err = m.Up()
	if err != nil {
		log.Fatalln("could not migrate up", err)
	}
}
