package config

import (
	"database/sql"
	"errors"
	"fmt"

	sqlite "github.com/glebarez/go-sqlite"
	dao "github.com/ryanfaerman/gophercnc/config/data"
)

type Resource struct {
	Name string
	Path string
}

func (c *config) AddResource(name, path, kind string) error {
	_, err := c.FindResourceByNameKind(name, kind)
	if err != sql.ErrNoRows {
		return errors.New("Resource exists")
	}

	_, err = c.queries.CreateResource(c.ctx, dao.CreateResourceParams{
		Name: name,
		Path: sql.NullString{String: path, Valid: true},
		Kind: sql.NullString{String: kind, Valid: true},
	})

	if err != nil {
		switch e := err.(type) {
		case *sqlite.Error:
			switch e.Code() {
			case 1555:
				return fmt.Errorf("resource already exists; %w", err)
			default:
				return fmt.Errorf("cannot add resource; %w", err)
			}
		}

		return fmt.Errorf("cannot add resource; %w", err)
	}

	active, err := c.ActiveResource(kind)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("cannot get active resource; %w", err)
	}

	if active.Name != name {
		if active.Name == "" {
			c.log.Debug("no active resources, activating", "name", name, "kind", kind)
			c.ActivateResource(name, kind)
		} else {
			c.log.Debug("another resource is active", "active", active.Name, "inactive", name, "kind", kind)
		}
	}

	return nil
}

func (c *config) FindResource() ([]Resource, error) {
	out := []Resource{}

	rs, err := c.queries.FindResource(c.ctx)
	if err != nil {
		return out, err
	}

	for _, r := range rs {
		out = append(out, Resource{
			Name: r.Name,
			Path: r.Path.String,
		})

	}

	return out, nil
}

func (c *config) FindResourceByKind(kind string) ([]Resource, error) {
	out := []Resource{}

	rs, err := c.queries.FindResourceByKind(c.ctx, sql.NullString{String: kind, Valid: true})
	if err != nil {
		return out, err
	}

	for _, r := range rs {
		out = append(out, Resource{
			Name: r.Name,
			Path: r.Path.String,
		})

	}

	return out, nil
}

func (c *config) FindResourceByNameKind(name, kind string) (Resource, error) {
	out := Resource{}

	r, err := c.queries.FindResourceByNameByKind(c.ctx, dao.FindResourceByNameByKindParams{
		Name: name,
		Kind: sql.NullString{String: kind, Valid: true},
	})
	if err != nil {
		return out, err
	}

	out.Name = r.Name
	out.Path = r.Path.String

	return out, nil
}

func (c *config) ActivateResource(name, kind string) error {
	_, err := c.FindResourceByNameKind(name, kind)
	if err != nil {
		return fmt.Errorf("Resource '%s' (%s) does not exist; %w", name, kind, err)
	}

	uri := fmt.Sprintf("%s.active", kind)

	return c.Set(uri, name)
}

func (c *config) ActiveResource(kind string) (Resource, error) {
	out := Resource{}

	uri := fmt.Sprintf("%s.active", kind)
	r, err := c.queries.ActiveResource(c.ctx, dao.ActiveResourceParams{
		Uri:  uri,
		Kind: sql.NullString{String: kind, Valid: true},
	})
	if err != nil {
		return out, err
	}

	out.Name = r.Name
	out.Path = r.Path.String

	return out, nil
}

func (c *config) RemoveResource(name, kind string) error {
	return c.queries.RemoveResource(c.ctx, dao.RemoveResourceParams{
		Name: name,
	})
}
