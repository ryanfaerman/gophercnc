package config

import (
	"database/sql"
	"errors"
	"fmt"

	sqlite "github.com/glebarez/go-sqlite"
	dao "github.com/ryanfaerman/gophercnc/config/data"
)

type Library struct {
	Name string
	Path string
}

// Libaries returns a map (name => file name) of every configured tool library.
// The files should exist within the library.path.
func (c *config) Libraries() ([]Library, error) {
	out := []Library{}
	libs, err := c.queries.ToolLibraries(c.ctx)
	if err != nil {
		return out, err
	}

	for _, lib := range libs {
		out = append(out, Library{
			Name: lib.Name,
			Path: lib.Path.String,
		})

	}

	return out, err
}

// GetLibraryPath returns the full path to the library for the given name.
func (c *config) GetLibrary(name string) (Library, error) {
	out := Library{}
	lib, err := c.queries.FindResourceByName(c.ctx, name)

	if err == nil {
		out.Name = lib.Name
		out.Path = lib.Path.String
	}

	if err != nil {
		err = fmt.Errorf("library '%s' not found; %w", name, err)
	}

	return out, err
}

// AddLibrary adds the given name and path to the library. The path should be a
// relative path from the library.path.
//
// - If the library already exists, an error will return
// - If there are no active libraries, it will be marked as active.
func (c *config) AddLibrary(name, path string) error {
	_, err := c.queries.FindResourceByName(c.ctx, name)
	if err != sql.ErrNoRows {
		return errors.New("Resource exists")
	}

	_, err = c.queries.AddToolLibrary(c.ctx, dao.AddToolLibraryParams{
		Name: name,
		Path: sql.NullString{String: path, Valid: true},
	})

	if err != nil {
		switch e := err.(type) {
		case *sqlite.Error:
			switch e.Code() {
			case 1555:
				return fmt.Errorf("library already exists; %w", err)
			default:
				return fmt.Errorf("cannot add library; %w", err)
			}
		}

		return fmt.Errorf("cannot add library; %w", err)
	}

	active, err := c.ActiveLibrary()
	if err != nil {
		return err
	}

	if active.Name != name {
		if active.Name == "" {
			c.log.Debug("no active libraries, activating", "name", name)
			c.ActivateLibrary(name)
		} else {
			c.log.Debug("another library is active", "active", active.Name, "inactive", name)
		}
	}

	return nil
}

// ActivateLibrary marks a tool library as the active library. The active
// library is intended for use in all toolpath generation operations.
func (c *config) ActivateLibrary(name string) error {
	if _, err := c.GetLibrary(name); err != nil {
		return fmt.Errorf("library '%s' does not exist; %w", name, err)
	}

	err := c.queries.SetActiveLibrary(c.ctx, sql.NullString{String: name, Valid: true})
	if err != nil {
		return fmt.Errorf("cannot activate library '%s'; %w", name, err)
	}
	return nil
}

// ActiveLibrary returns the active library
func (c *config) ActiveLibrary() (Library, error) {
	out := Library{}
	active, err := c.queries.ActiveToolLibrary(c.ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return out, nil
		}
		return out, fmt.Errorf("cannot get active library; %w", err)
	}

	out.Name = active.Name
	out.Path = active.Path.String

	return out, nil
}

// RemoveLibrary deletes the resource at the given name. If it is an activated
// library, the underlying SQL JOIN will return an empty Library, indicating
// that there are no active libraries. Activating a different tool library is
// left as an exercise for the consumer.
func (c *config) RemoveLibrary(name string) error {
	return c.queries.RemoveResource(c.ctx, name)
}
