package sqlite_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

func TestMain(m *testing.M) {
	goose.SetDialect("sqlite3")
	code := m.Run()
	os.Exit(code)
}

func setupSqlite(name string) (db *sqlx.DB) {
	connectionString := fmt.Sprintf("file:%s?mode=memory&cache=shared&_pragma=foreign_keys=ON", name)
	db = sqlx.MustConnect("sqlite", connectionString)

	err := goose.Up(db.DB, "../../../db/migrations")
	if err != nil {
		panic(fmt.Errorf("Migrating the database failed with: %w", err))
	}
	err = goose.Up(db.DB, "../../../db/seeds")
	if err != nil {
		panic(fmt.Errorf("Seeding the database failed with: %w", err))
	}

	return
}
