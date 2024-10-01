package orm_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/bbfh-dev/go-orm/orm"
	"github.com/bbfh-dev/go-orm/orm/tables"
	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/go-libsql"
	"gotest.tools/assert"
)

type TestTableA struct {
	ExampleId   int64  `db:"example_id"  create:"INTEGER PRIMARY KEY"`
	Bool        bool   `db:"bool"        create:"BOOLEAN NOT NULL DEFAULT 0"`
	Name        string `db:"name"        create:"TEXT NOT NULL DEFAULT ''"`
	Description string `db:"description" create:"TEXT NOT NULL DEFAULT ''"`
}

type TestTableB struct {
	ExampleId int64 `db:"example_id" create:"INTEGER PRIMARY KEY"`
	Field     int   `db:"field"      create:"INTEGER"`
}

type TestTableC struct {
	ExampleId int64 `db:"example_id" create:"INTEGER PRIMARY KEY"`
	Field     int   `db:"field"      create:"INTEGER NOT NULL DEFAULT 0"`
}

type TestTableD struct {
	ExampleId int64  `db:"example_id" create:"INTEGER PRIMARY KEY"`
	Field     int    `db:"field"      create:"INTEGER"`
	Name      string `db:"name"       create:"TEXT NOT NULL DEFAULT ''"`
}

type TestTableE struct {
	ExampleId int64 `db:"example_id" create:"INTEGER PRIMARY KEY"`
}

func (table TestTableA) SQL() string {
	return "example_1"
}

func (table TestTableB) SQL() string {
	return "example_2"
}

func (table TestTableC) SQL() string {
	return "example_2"
}

func (table TestTableD) SQL() string {
	return "example_2"
}

func (table TestTableE) SQL() string {
	return "example_2"
}

func TestDB(test *testing.T) {
	// This test is a fucking abomination, but it works! SO DON'T BREAK IT

	slog.SetLogLoggerLevel(slog.LevelDebug)
	os.Setenv(orm.MIGRATE_ENV, "1")
	os.Setenv(orm.INSTANT_ENV, "1")

	path := filepath.Join(os.TempDir(), "test.db")
	test.Cleanup(func() {
		os.Remove(path)
	})

	sqlDB, err := sqlx.Open("libsql", "file:"+path)
	assert.NilError(test, err)
	assert.Assert(test, sqlDB != nil)

	db := orm.NewDB(sqlDB)
	assert.Assert(test, db != nil)

	db.Tables = []tables.Table{TestTableA{}, TestTableB{}}
	assert.NilError(test, db.GenMigrations())

	db.Tables = []tables.Table{TestTableA{}, TestTableC{}}
	assert.NilError(test, db.GenMigrations())

	db.Tables = []tables.Table{TestTableA{}, TestTableD{}}
	assert.NilError(test, db.GenMigrations())

	pragma, err := db.PragmaOf(TestTableB{})
	assert.NilError(test, err)

	assert.DeepEqual(test, map[string]string(pragma), map[string]string{
		"example_id": "INTEGER PRIMARY KEY",
		"field":      "INTEGER",
		"name":       "TEXT NOT NULL DEFAULT ''",
	})

	db.Tables = []tables.Table{TestTableA{}, TestTableE{}}
	assert.NilError(test, db.GenMigrations())

	pragma, err = db.PragmaOf(TestTableB{})
	assert.NilError(test, err)

	assert.DeepEqual(test, map[string]string(pragma), map[string]string{
		"example_id": "INTEGER PRIMARY KEY",
	})

	pragma, err = db.PragmaOf(TestTableA{})
	assert.NilError(test, err)

	assert.DeepEqual(test, map[string]string(pragma), map[string]string{
		"example_id":  "INTEGER PRIMARY KEY",
		"bool":        "BOOLEAN NOT NULL DEFAULT 0",
		"name":        "TEXT NOT NULL DEFAULT ''",
		"description": "TEXT NOT NULL DEFAULT ''",
	})

	err = db.InsertEntity(TestTableA{
		ExampleId:   1,
		Bool:        true,
		Name:        "Hello",
		Description: "World!",
	})
	assert.NilError(test, err)

	example := TestTableA{
		ExampleId:   1,
		Bool:        true,
		Name:        "Hello",
		Description: "World!",
	}

	var tables []TestTableA
	err = orm.GetEntities(db, &tables, "")
	assert.NilError(test, err)
	assert.DeepEqual(test, tables, []TestTableA{example})

	var table TestTableA
	err = orm.GetSingleEntity(db, &table, `WHERE example_id = 1`)
	assert.NilError(test, err)
	assert.DeepEqual(test, table, example)

	fields, err := db.Fields("", TestTableA{}, "name", "description")
	assert.NilError(test, err)
	assert.DeepEqual(test, map[string]string(fields), map[string]string{
		"name":        "Hello",
		"description": "World!",
	})

	_, err = db.Fields("", TestTableA{}, "unknown")
	assert.Check(test, err != nil)
}
