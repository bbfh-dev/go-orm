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

// func GetFields[T any](db *DB, table tables.Table, out *[]T, fields []string, query string) error {
// 	rows, err := db.handle.Queryx(
// 		fmt.Sprintf("SELECT %s FROM %s %s;", strings.Join(fields, ", "), table.SQL(), query),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()
//
// 	// cols, err := rows.Columns()
// 	// if err != nil {
// 	// 	return err
// 	// }
//
// 	for rows.Next() {
// 		row := make(map[string]interface{})
// 		err := rows.MapScan(row)
// 		if err != nil {
// 			return tools.PrefixErr("Scanning row", err)
// 		}
// 		for i, field := range fields {
// 			value, ok := row[field]
// 			if !ok {
// 				return fmt.Errorf("SQL didn't return field %s", field)
// 			}
// 			*out[i] = value
// 		}
// 		fmt.Println(row)
// 	}
//
// 	return nil
// }
