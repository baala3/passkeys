package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/migrate"
)

func GetTestDB() *bun.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	testDB := bun.NewDB(db, sqlitedialect.New())

	// run migrations
	ctx := context.Background()
	migrations := migrate.NewMigrations()

	if err = migrations.DiscoverCaller(); err != nil {
		panic(err)
	}

	migrator := migrate.NewMigrator(testDB, migrations)

	if err = migrator.Init(ctx); err != nil {
		panic(err)
	}

	if _, err = migrator.Migrate(ctx); err != nil {
		fmt.Println("error migrating test db", err)
		panic(err)
	}
	return testDB
}
