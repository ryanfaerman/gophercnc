// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package data

import (
	"context"
	"database/sql"
)

const activeToolLibrary = `-- name: ActiveToolLibrary :one
SELECT resources.name, resources.path FROM configs 
JOIN resources ON configs.data = resources.name 
WHERE configs.uri = "library.active"
`

type ActiveToolLibraryRow struct {
	Name string
	Path sql.NullString
}

func (q *Queries) ActiveToolLibrary(ctx context.Context) (ActiveToolLibraryRow, error) {
	row := q.db.QueryRowContext(ctx, activeToolLibrary)
	var i ActiveToolLibraryRow
	err := row.Scan(&i.Name, &i.Path)
	return i, err
}

const addToolLibrary = `-- name: AddToolLibrary :one
INSERT INTO resources (name, path, kind)
VALUES(?1, ?2, "tool.library")
RETURNING name, path, format, kind
`

type AddToolLibraryParams struct {
	Name string
	Path sql.NullString
}

func (q *Queries) AddToolLibrary(ctx context.Context, arg AddToolLibraryParams) (Resource, error) {
	row := q.db.QueryRowContext(ctx, addToolLibrary, arg.Name, arg.Path)
	var i Resource
	err := row.Scan(
		&i.Name,
		&i.Path,
		&i.Format,
		&i.Kind,
	)
	return i, err
}

const findResourceByName = `-- name: FindResourceByName :one
SELECT name, path, format, kind FROM resources
WHERE name = ?
`

func (q *Queries) FindResourceByName(ctx context.Context, name string) (Resource, error) {
	row := q.db.QueryRowContext(ctx, findResourceByName, name)
	var i Resource
	err := row.Scan(
		&i.Name,
		&i.Path,
		&i.Format,
		&i.Kind,
	)
	return i, err
}

const getConfig = `-- name: GetConfig :one
SELECT data FROM configs 
WHERE uri = ?
`

func (q *Queries) GetConfig(ctx context.Context, uri string) (sql.NullString, error) {
	row := q.db.QueryRowContext(ctx, getConfig, uri)
	var data sql.NullString
	err := row.Scan(&data)
	return data, err
}

const removeResource = `-- name: RemoveResource :exec
DELETE FROM resources WHERE name=?
`

func (q *Queries) RemoveResource(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, removeResource, name)
	return err
}

const setActiveLibrary = `-- name: SetActiveLibrary :exec
INSERT INTO configs (uri, data) 
VALUES("library.active", ?1)
ON CONFLICT(uri) DO UPDATE 
SET data = ?1
`

func (q *Queries) SetActiveLibrary(ctx context.Context, data sql.NullString) error {
	_, err := q.db.ExecContext(ctx, setActiveLibrary, data)
	return err
}

const setConfig = `-- name: SetConfig :exec
INSERT INTO configs (uri, data)
VALUES (?1, ?2)
ON CONFLICT(uri) DO UPDATE
SET data = ?2
`

type SetConfigParams struct {
	Uri  string
	Data sql.NullString
}

func (q *Queries) SetConfig(ctx context.Context, arg SetConfigParams) error {
	_, err := q.db.ExecContext(ctx, setConfig, arg.Uri, arg.Data)
	return err
}

const toolLibraries = `-- name: ToolLibraries :many
SELECT name, path, format, kind FROM resources
WHERE kind = "tool.library"
`

func (q *Queries) ToolLibraries(ctx context.Context) ([]Resource, error) {
	rows, err := q.db.QueryContext(ctx, toolLibraries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Resource
	for rows.Next() {
		var i Resource
		if err := rows.Scan(
			&i.Name,
			&i.Path,
			&i.Format,
			&i.Kind,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}