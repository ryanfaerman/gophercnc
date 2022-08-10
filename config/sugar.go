package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func Load() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(configDir, ApplicationName)
	if err := os.MkdirAll(path, 0750); err != nil {
		return err
	}

	return c.Load(path)
}
func Reset() error { return c.Reset() }
func Set(uri, data string) error {
	if !c.loaded {
		return errors.New("config is not loaded")
	}

	for _, p := range protected {
		if p == uri {
			return fmt.Errorf("uri '%s' is reserved/protected", uri)
		}
	}

	return c.Set(uri, data)
}
func Get(uri string) (string, error) {
	if !c.loaded {
		return "", errors.New("config is not loaded")
	}
	return c.Get(uri)
}

func GetAll() ([]ConfigOption, error) {

	if !c.loaded {
		return []ConfigOption{}, errors.New("config is not loaded")
	}
	return c.GetAll()
}

func Unset(uri string) error {
	if !c.loaded {
		return errors.New("config is not loaded")
	}
	for _, p := range protected {
		if p == uri {
			return fmt.Errorf("uri '%s' is reserved/protected", uri)
		}
	}
	return c.Unset(uri)
}

func Close() error { return c.Close() }

// Library-related Sweets

func LibraryCachePath() string {
	p, err := c.Get("library.path")
	if err != nil {
		panic(err.Error())
	}
	return p
}

// func Libraries() ([]Library, error)           { return c.Libraries() }
// func GetLibrary(name string) (Library, error) { return c.GetLibrary(name) }
// func AddLibrary(name, path string) error { return c.AddLibrary(name, path) }
// func ActivateLibrary(name string) error { return c.ActivateLibrary(name) }
// func ActiveLibrary() (Library, error) { return c.ActiveLibrary() }
// func RemoveLibrary(name string) error { return c.RemoveLibrary(name) }

func Libraries() ([]Resource, error)           { return c.FindResourceByKind("tool.library") }
func GetLibrary(name string) (Resource, error) { return c.FindResourceByNameKind(name, "tool.library") }
func AddLibrary(name, path string) error       { return c.AddResource(name, path, "tool.library") }
func ActivateLibrary(name string) error        { return c.ActivateResource(name, "tool.library") }
func ActiveLibrary() (Resource, error)         { return c.ActiveResource("tool.library") }
func RemoveLibrary(name string) error          { return c.RemoveResource(name, "tool.library") }

// Machine-related Sweets
func MachineCachePath() string {
	p, err := c.Get("machine.path")
	if err != nil {
		panic(err.Error())
	}
	return p
}
func Machines() ([]Resource, error)            { return c.FindResourceByKind("machine") }
func GetMachine(name string) (Resource, error) { return c.FindResourceByNameKind(name, "machine") }
func AddMachine(name, path string) error       { return c.AddResource(name, path, "machine") }
func ActivateMachine(name string) error        { return c.ActivateResource(name, "machine") }
func ActiveMachine() (Resource, error)         { return c.ActiveResource("machine") }
func RemoveMachine(name string) error          { return c.RemoveResource(name, "machine") }
