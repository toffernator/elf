package sqlite_test

import (
	"fmt"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

var db = sqlx.MustConnect("sqlite", ":memory:")

var didSeedRun bool = false

func seed() {
	if didSeedRun {
		return
	}

	goose.SetDialect("sqlite3")
	err := goose.Up(db.DB, "../../db/migrations")
	if err != nil {
		panic(fmt.Errorf("Migrating the database failed with: %w", err))
	}
	err = goose.Up(db.DB, "../../db/seeds")
	if err != nil {
		panic(fmt.Errorf("Seeding the database failed with: %w", err))
	}

	didSeedRun = true
}
