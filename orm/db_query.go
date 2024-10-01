package orm

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bbfh-dev/go-orm/orm/tables"
	"github.com/bbfh-dev/go-tools/tools/terr"
)

func (db *DB) Select(dest interface{}, query string) error {
	terr.Assert(dest != nil, "Columns must not be nil")
	terr.Assert(db.handle != nil, "DB handle must not be nil")
	slog.Debug("(ORM) Query", "query", query)

	err := db.handle.Select(dest, query)
	return terr.Prefix("(ORM) Query "+query, err)
}

func (db *DB) Exec(query string, args ...any) error {
	query = fmt.Sprintf(query, args...)
	slog.Debug("(ORM) Exec", "query", query)

	_, err := db.handle.Exec(query)
	return err
}

type fields map[string]string

func (db *DB) Fields(query string, table tables.Table, fields ...string) (fields, error) {
	var out = map[string]string{}

	rows, err := db.handle.Queryx(fmt.Sprintf(
		"SELECT %s FROM %s %s;",
		strings.Join(fields, ", "),
		table.SQL(),
		query,
	))
	if err != nil {
		return out, terr.Prefix("DB Query(SELECT)", err)
	}
	defer rows.Close()

	for rows.Next() {
		row := make(map[string]interface{})
		err = rows.MapScan(row)
		if err != nil {
			return out, terr.Prefix("DB Rows.Scan", err)
		}
		for key, value := range row {
			out[key] = fmt.Sprintf("%v", value)
		}
	}

	return out, nil
}
