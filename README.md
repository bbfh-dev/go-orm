# Go ORM

Designed to be a simple and reliable ORM for SQL databases.

It doesn't perform any database altering operations unless given permission to. It's explicit in its actions and prints useful debug information using `slog.Debug`.

<!-- vim-markdown-toc GFM -->

* [Features](#features)
* [Example](#example)

<!-- vim-markdown-toc -->

# Features

- Automatically syncronize database table with Go struct (CREATE, ALTER, DROP);
- Query fields / entire tables with ease;
- Does not perform any operations unless ORM_MIGRATE=1 env is present;
- Explicit errors prefixed with the operation that caused them;
- Explicit logging using `slog` (With DEBUG messages per every SQL call)

# Example

```go
// An example struct, every field must have "db" and "create" tags
type TestTable struct {
	ExampleId   int64  `db:"example_id"  create:"INTEGER PRIMARY KEY"`
	Name        string `db:"name"        create:"TEXT NOT NULL DEFAULT ''"`
	Description string `db:"description" create:"TEXT NOT NULL DEFAULT ''"`
}

// The only required method, returns the name of the table in SQL database
func (table TestTable) SQL() string {
	return "test_table"
}

func main() {
	// In this example we use SQLite3
	sqlDB, err := sqlx.Open("libsql", "file:/tmp/my.db")
	if err != nil {
		panic(err)
	}

	// Create the Go-ORM database wrapper
	db := orm.NewDB(sqlDB)
	// Tell it what tables we have
	db.Tables = []tables.Table{TestTable{}}

	// Ensure that SQL database is in-sync with our go struct
	err = db.GenMigrations()
	if err != nil {
		panic(err)
	}
}
```
