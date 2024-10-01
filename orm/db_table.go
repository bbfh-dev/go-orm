package orm

import (
	"errors"
	"fmt"

	"github.com/bbfh-dev/go-orm/orm/tables"
)

var EmptyErr = errors.New("Result is empty")

func (db *DB) InsertEntity(table tables.Table) error {
	tableValues := tables.Values(table)

	var keys, values []string
	for key, value := range tableValues {
		keys = append(keys, key)
		values = append(values, value)
	}

	return db.Exec(tables.INSERT_VALUES(
		table,
		keys,
		values,
	))
}

func Entities[T tables.Table](db *DB, table *[]T, query string) error {
	var in = *new(T)
	return db.Select(table, fmt.Sprintf("SELECT * FROM %s %s;", in.SQL(), query))
}

func SingleEntity[T tables.Table](db *DB, table *T, query string) error {
	var entities []T
	err := Entities(db, &entities, query)
	if err != nil {
		return err
	}

	if len(entities) < 1 {
		return EmptyErr
	}

	*table = entities[0]
	return nil
}
