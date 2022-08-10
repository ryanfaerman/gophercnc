package machine

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ryanfaerman/gophercnc/config"
	"github.com/segmentio/ksuid"
	"gopkg.in/yaml.v2"
)

type Range struct {
	Minimum float64 `json:"minimum" yaml:"minimum"`
	Maximum float64 `json:"maximum" yaml:"maximum"`
}

type Axis struct {
	ID    byte  `json:"id" yaml:"id"`
	Range Range `json:"range" yaml:"range"`
}

type Capabilities struct {
	Additive   bool `json:"additive" yaml:"additive"`
	Inspection bool `json:"inspection" yaml:"inspection"`
	Jet        bool `json:"jet" yaml:"jet"`
	Milling    bool `json:"milling" yaml:"milling"`
	Turning    bool `json:"turning" yaml:"turning"`
}

type Machine struct {
	Name         string            `json:"name" yaml:"name"`
	Axes         []Axis            `json:"axes" yaml:"axes"`
	Capabilities Capabilities      `json:"capabilities" yaml:"capabilities"`
	Spindle      Range             `json:"spindle" yaml:"spindle"`
	Procedures   map[string]string `json:"procedures" yaml:"procedures"`
}

func LoadMachine(path string) (Machine, error) {
	// TODO: load *.machine files

	var m Machine

	switch filepath.Ext(path) {
	case ".machine":
		return m, errors.New("currently unsupported machine kind")
	case ".yml", ".yaml":
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return m, err
		}

		err = yaml.Unmarshal(data, &m)

		return m, err

	}
	return m, nil
}

func Import(path string) error {
	l := config.Logger.WithFields("fn", "machine.Import")

	if err := CleanMachineCache(); err != nil {
		return err
	}

	l.Debug("importing machine definition", "path", path)
	m, err := LoadMachine(path)
	if err != nil {
		return err
	}

	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	targetName := ksuid.New().String() + filepath.Ext(path)
	targetPath := filepath.Join(config.MachineCachePath(), targetName)

	if err := os.MkdirAll(config.MachineCachePath(), 0750); err != nil {
		return err
	}

	dst, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	l.Debug("copying machine to cache", "name", m.Name, "src", path, "dst", targetPath)
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	l.Debug("adding library to config", "name", m.Name)
	if err := config.AddMachine(m.Name, targetPath); err != nil {

		return fmt.Errorf("Cannot import library; %w", err)
	}
	return nil
}

// CleanLibraryCache compares the libraries on disk to the ones we have
// configured. It then removes any on disk from the cache that aren't
// present in the config.
//
// If there are any libaries in the config that don't exist on disk, they are
// removed from the config.
func CleanMachineCache() error {
	l := config.Logger.WithFields("fn", "machine.CleanLibraryCache")

	l.Debug("loading stored machines")
	machines, err := config.Machines()
	if err != nil {
		return err
	}

	validPaths := make(map[string]string)

	for _, m := range machines {
		validPaths[filepath.Base(m.Path)] = m.Name

		if _, err := os.Stat(m.Path); os.IsNotExist(err) {
			l.Warn("removing invalid machine from the config", "name", m.Name, "path", m.Path)
			if err := config.RemoveMachine(m.Name); err != nil {
				l.WithError(err).Error("cannot remove invalid machine", "name", m.Name)
			}

		}
	}

	l.Debug("Reading machine cache", "path", config.MachineCachePath())
	files, err := ioutil.ReadDir(config.MachineCachePath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, f := range files {
		if _, ok := validPaths[f.Name()]; !ok {
			path := filepath.Join(config.MachineCachePath(), f.Name())
			l.Warn("removing orphan machine", "path", path)

			os.Remove(path)
		}
	}

	return nil
}
