package postgres

import (
	"database/sql"
	"gopher/internal/config"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(config config.Config) {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("pgx", config.DatabaseURL())
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}
	db.Close()
}
