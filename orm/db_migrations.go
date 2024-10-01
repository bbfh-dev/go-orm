package orm

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bbfh-dev/go-orm/orm/tables"
	"github.com/bbfh-dev/go-tools/tools/terr"
	"github.com/bbfh-dev/go-tools/tools/tlog"
)

const MIGRATE_ENV = "ORM_MIGRATE"
const INSTANT_ENV = "ORM_INSTANT"

const COUNTDOWN = 5

func (db *DB) GenMigrations() error {
	if len(db.Tables) == 0 {
		slog.Warn("(ORM) DB.Tables is empty! Are you sure you set it up correctly?")
		return nil
	}

	slog.Info("--- (ORM) Performing Database migration generation")

	var altered = false
	var migrations []string

	for _, table := range db.Tables {
		pragmaColumns, err := db.PragmaOf(table)
		if err != nil {
			return err
		}

		if IsPragmaEmpty(pragmaColumns) {
			slog.Warn("Creating table", "name", table.SQL())
			migrations = append(migrations, tables.CREATE_TABLE(table))
			continue
		}

		var tableColumns = tables.GetColumns(table)

		for column, tableCreate := range tableColumns {
			pragmaCreate, ok := pragmaColumns[column]
			if !ok {
				migrations = append(migrations, tables.ALTER_TABLE_ADD(table, column, tableCreate))
				continue
			}
			if tableCreate == pragmaCreate {
				continue
			}
			altered = true
			migrations = append(migrations, tables.CREATE_TEMP_TABLE(table))
			migrations = append(
				migrations,
				tables.COPY_TABLE(table, table.SQL()+"__tmp", pragmaColumns),
			)
			migrations = append(migrations, tables.DROP_TABLE(table))
			migrations = append(
				migrations,
				tables.ALTER_TABLE_RENAME(table.SQL()+"__tmp", table.SQL()),
			)
		}

		for column := range pragmaColumns {
			_, ok := tableColumns[column]
			if !ok {
				migrations = append(migrations, tables.ALTER_TABLE_DROP(table, column))
			}
		}
	}

	slog.Info("=== (ORM) Finished Database migration generation")

	if len(migrations) != 0 {
		env, ok := os.LookupEnv(MIGRATE_ENV)
		if !ok || env == "0" {
			tlog.Warn(
				"(ORM) You have %d unapplied migrations, meaning that your database is out of sync! Please create a backup and provide %s=1 environment variable to apply them.",
				len(migrations),
				MIGRATE_ENV,
			)
			tlog.Info("(ORM) List of migrations to apply:")
			for _, migration := range migrations {
				slog.Info(migration)
			}
		} else {
			return terr.Prefix("Database Migration", db.ApplyMigrations(migrations, altered))
		}
	}

	return nil
}

func (db *DB) ApplyMigrations(migrations []string, altered bool) error {
	env, ok := os.LookupEnv(INSTANT_ENV)
	if ok && env != "0" {
		slog.Info("Performing instant migration (no safety countdown)")
		return db.ApplyMigrationsNow(migrations, altered)
	}

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	for i := COUNTDOWN; i > 0; i-- {
		tlog.Warn("(DANGER!) Migrating in %d...\n", i)
		time.Sleep(1 * time.Second)

		select {
		case <-interruptChan:
			slog.Info("Migration interrupted by the user.")
			return nil
		default:
			// Continue with countdown
		}
	}

	slog.Info("Running database migration now!")
	return db.ApplyMigrationsNow(migrations, altered)
}

func (db *DB) ApplyMigrationsNow(migrations []string, altered bool) error {
	length := len(migrations)

	for i, migration := range migrations {
		if altered && strings.Contains(migration, "\nADD") {
			tlog.Info("--- Skipping migration %d/%d (Overwritten)", i+1, length)
			continue
		}

		tlog.Info("--- Applying migration %d/%d...", i+1, length)
		err := db.Exec(migration)
		if err != nil {
			slog.Error(err.Error())
		}
	}

	return nil
}
