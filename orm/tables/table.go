package tables

import (
	"fmt"
	"reflect"

	"github.com/bbfh-dev/go-tools/tools/terr"
)

type Table interface {
	// The SQL name of the table
	SQL() string
}

func assertFieldDBExists(field reflect.StructField) {
	terr.Assert(
		len(field.Tag.Get("db")) > 0,
		"Table field must have 'db' tag describing how it should be called in SQL!",
	)
}

func Columns(table Table) map[string]string {
	tableType := reflect.TypeOf(table)
	var out = map[string]string{}

	for i := range tableType.NumField() {
		field := tableType.Field(i)
		assertFieldDBExists(field)
		terr.Assert(
			len(field.Tag.Get("create")) > 0,
			"Table field must have 'create' tag describing the datatype and constaints (e.g. TEXT NOT NULL)!",
		)
		out[field.Tag.Get("db")] = field.Tag.Get("create")
	}

	return out
}

func Values(table Table) map[string]string {
	tableType := reflect.TypeOf(table)
	tableValue := reflect.ValueOf(table)
	var out = map[string]string{}

	for i := range tableType.NumField() {
		typeField := tableType.Field(i)
		valueField := tableValue.Field(i)

		assertFieldDBExists(typeField)
		switch valueField.Kind() {
		case reflect.String:
			out[typeField.Tag.Get("db")] = fmt.Sprintf("'%s'", valueField.String())
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Float32, reflect.Float64:
			out[typeField.Tag.Get("db")] = fmt.Sprintf("%d", valueField.Int())
		default:
			out[typeField.Tag.Get("db")] = fmt.Sprintf("%v", valueField)
		}
	}

	return out
}
