// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package data

import (
	"database/sql"
)

type Config struct {
	Uri  string
	Data sql.NullString
}

type Resource struct {
	Name   string
	Path   sql.NullString
	Format sql.NullString
	Kind   sql.NullString
}