package config

import (
	"fmt"
	"sort"
)

type ConfigOption struct {
	Uri  string
	Data string
}

func (c *config) GetAll() ([]ConfigOption, error) {
	out := []ConfigOption{}

	opts, err := c.queries.Configs(c.ctx)
	if err != nil {
		return out, fmt.Errorf("cannot get configs; %w", err)
	}

	for _, opt := range opts {
		out = append(out, ConfigOption{
			Uri:  opt.Uri,
			Data: opt.Data.String,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Uri < out[j].Uri
	})

	return out, nil
}

func (c *config) Unset(uri string) error {
	return c.queries.UnsetConfig(c.ctx, uri)
}
