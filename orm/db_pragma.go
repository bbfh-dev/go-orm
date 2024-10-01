package orm

import (
	"fmt"
	"strings"

	"github.com/bbfh-dev/go-orm/orm/tables"
	"github.com/bbfh-dev/go-tools/tools"
)

type pragmaCol struct {
	CID        int     `db:"cid"`
	Name       string  `db:"name"`
	Type       string  `db:"type"`
	NotNull    int     `db:"notnull"`
	DefaultVal *string `db:"dflt_value"`
	PrimaryKey int     `db:"pk"`
}

type pragma map[string]string

func (db *DB) PragmaOf(table tables.Table) (pragma, error) {
	var columns []pragmaCol
	err := db.Select(&columns, fmt.Sprintf("PRAGMA table_info(%s);", table.SQL()))
	if err != nil {
		return nil, tools.PrefixErr("PragmaOf()", err)
	}

	var builder strings.Builder
	var out = map[string]string{}

	for _, col := range columns {
		builder.WriteString(col.Type)
		if col.PrimaryKey == 1 {
			builder.WriteString(" PRIMARY KEY")
		}
		if col.NotNull == 1 {
			builder.WriteString(" NOT NULL")
		}
		if col.DefaultVal != nil {
			builder.WriteString(" DEFAULT ")
			builder.WriteString(*col.DefaultVal)
		}
		out[col.Name] = builder.String()
		builder.Reset()
	}

	return out, nil
}

func IsPragmaEmpty(rows pragma) bool {
	return len(rows) == 0
}
