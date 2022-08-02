package config

import (
	"context"
	"database/sql"
	"embed"
	"os"
	"path/filepath"
	"sync"

	goose "github.com/pressly/goose/v3"
	dao "github.com/ryanfaerman/gophercnc/config/data"
	"github.com/ryanfaerman/gophercnc/log"
)

var (
	ApplicationName = "gophercnc"
	defaults        = map[string]string{
		"cache.path":   "",
		"library.path": "",
	}
	protected = []string{
		"cache.path",
		"library.path",
	}

	c      config
	Logger log.Logger
)

//go:embed migrations/*.sql
var migrations embed.FS

// There are some values in the Defaults that we cannot determine before
// runtime. This gives us a hook to do that.
func init() {
	cache, err := os.UserCacheDir()
	if err != nil {
		panic(err.Error())
	}
	defaults["cache.path"] = filepath.Join(cache, "gophercnc")
	defaults["library.path"] = filepath.Join(defaults["cache.path"], "libraries")

}

type config struct {
	path    string
	ctx     context.Context
	db      *sql.DB
	queries *dao.Queries
	once    sync.Once
	loaded  bool
	log     log.Logger
}

// Load the config system, storing data at the given path. This can safely be
// called multiple times.
//
// The underlying config data is stored in SQLite. During the Load, all
// migrations are applied and the queries DAO is setup as well. This should
// ensure that the config database has all the right schema. Migrations are
// embedded into the resulting binary and should be available even without the
// source.
func (c *config) Load(path string) error {
	var err error
	c.log = Logger.WithFields("pkg", "config")

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		path = filepath.Join(path, "config.db")
	}
	c.path = path

	if c.path != path {
		c.once = sync.Once{}
	}

	c.once.Do(func() {
		Logger.Debug("loading config", "path", c.path)

		c.ctx = context.Background()

		c.db, err = sql.Open("sqlite", c.path)

		// Setup our migrations here
		goose.SetLogger(Logger)
		goose.SetBaseFS(migrations)

		err = goose.SetDialect("sqlite")
		if err != nil {
			return
		}

		err = goose.Up(c.db, "migrations")
		if err != nil {
			return
		}

		c.queries = dao.New(c.db)

		err = c.SetDefaults()

		if err == nil {
			c.loaded = true
		}

	})

	return err
}

// SetDefaults ensures that certain data exists in the config at all times.
// These are poorly named, but they will reset on every load. When changing the
// values, they must be changed after every Load.
func (c *config) SetDefaults() error {
	for uri, data := range defaults {
		c.log.Debug("setting default", "uri", uri, "data", data)
		if err := c.Set(uri, data); err != nil {
			return err
		}
	}

	return nil
}

// Reset the configuration DB entirely. This is currently done by just deleting
// the underlying SQLite database.
func (c *config) Reset() error {
	c.log.Debug("resetting the config db")

	return os.Remove(c.path)
}

// Set the given config URI to the given value. This will replace the value if
// it is already defined. The database is assumed to be available by others
// calling Load.
func (c *config) Set(uri string, data string) error {
	return c.queries.SetConfig(c.ctx, dao.SetConfigParams{
		Uri:  uri,
		Data: sql.NullString{String: data, Valid: data != ""},
	})
}

// Get the data for the given uri. If undefined, an empty string is returned.
// An error should only be returned if there is some underlying system problem.
// The database is assumed to be available by others calling Load.
func (c *config) Get(uri string) (string, error) {
	v, err := c.queries.GetConfig(c.ctx, uri)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	if v.Valid {
		return v.String, nil
	}

	return "", nil
}

// Close the config and allow any cleanup to occur as part of shutdown. Once
// closed, the config cannot be used.
func (c *config) Close() error {
	return c.db.Close()
}
