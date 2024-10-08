package tables

import (
	"fmt"
	"strings"

	"github.com/bbfh-dev/go-tools/tools/tmap"
)

func createTable(table Table, name string) string {
	return fmt.Sprintf(
		"CREATE TABLE %s (%s);",
		name,
		strings.Join(
			tmap.Flatten(
				Columns(table),
				func(k, v string) string { return fmt.Sprintf("'%s' %s", k, v) },
			),
			", ",
		),
	)
}

func CREATE_TABLE(table Table) string {
	return createTable(table, table.SQL())
}

func CREATE_TEMP_TABLE(table Table) string {
	return createTable(table, table.SQL()+"__tmp")
}

func COPY_TABLE(table Table, dest string, pragma map[string]string) string {
	keys, _ := tmap.Split(pragma)
	fields := strings.Join(keys, ", ")

	return fmt.Sprintf(
		`INSERT INTO %s (%s)
		SELECT %s FROM %s;`,
		dest,
		fields,
		fields,
		table.SQL(),
	)
}

func DROP_TABLE(table Table) string {
	return fmt.Sprintf("DROP TABLE %s;", table.SQL())
}

func ALTER_TABLE_RENAME(oldname string, name string) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", oldname, name)
}

func ALTER_TABLE_ADD(table Table, column string, create string) string {
	return fmt.Sprintf("ALTER TABLE %s\nADD %s %s", table.SQL(), column, create)
}

func ALTER_TABLE_DROP(table Table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s\nDROP COLUMN %s", table.SQL(), column)
}

func INSERT_VALUES(table Table, fields []string, values []string) string {
	return fmt.Sprintf(
		"INSERT INTO %s (%s)\nVALUES (%s);",
		table.SQL(),
		strings.Join(fields, ", "),
		strings.Join(values, ", "),
	)
}
