// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package tutorial

import (
	"database/sql"
)

type Author struct {
	ID   int64
	Name string
	Bio  sql.NullString
}

type Library struct {
	Name   string
	Path   string
	Active sql.NullBool
}
