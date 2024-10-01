package orm

import (
	"github.com/bbfh-dev/go-orm/orm/tables"
	"github.com/jmoiron/sqlx"
)

// Database handle wrapper.
type DB struct {
	handle *sqlx.DB
	// The slice of registered tables.
	//
	// Modify it DIRECTLY, because why build useless interfaces.
	Tables []tables.Table
}

// Create a new DB object using the handle.
//
// Note: Don't forget to add your table.Table-s into DB.Tables after creation!
//
// In order to use the migration functionality, call DB.GenMigrations()
func NewDB(db *sqlx.DB) *DB {
	return &DB{
		handle: db,
	}
}

// Access the underlying *sqlx.DB directly
func (db *DB) Handle() *sqlx.DB {
	return db.handle
}
