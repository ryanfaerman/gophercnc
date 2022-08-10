package tool

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ryanfaerman/gophercnc/config"
	"github.com/segmentio/ksuid"
)

type Library struct {
	Tools   []Tool `json:"data"`
	Version int    `json:"version"`
}

func (l Library) FindByNumber(n int) (Tool, error) {
	for _, t := range l.Tools {
		if t.Number() == n {
			return t, nil
		}
	}
	return Tool{}, errors.New("Tool Not Found")
}

func IsValidLibrary(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

func LoadLibrary(path string) (Library, error) {

	if !IsValidLibrary(path) {
		return Library{}, errors.New("invalid library")
	}

	switch filepath.Ext(path) {
	case ".tools":
		r, err := zip.OpenReader(path)
		if err != nil {
			return Library{}, errors.New("Cannot open path")
		}
		defer r.Close()

		validLibraryFile := false
		for _, f := range r.File {
			if f.Name != "tools.json" {
				continue
			}
			validLibraryFile = true

			tf, err := f.Open()
			if err != nil {
				return Library{}, err
			}
			defer tf.Close()

			data, err := ioutil.ReadAll(tf)
			if err != nil {
				return Library{}, err
			}

			var lib Library
			return lib, json.Unmarshal(data, &lib)
		}

		if !validLibraryFile {
			return Library{}, errors.New("invalid library file, missing tools.json")
		}
	case ".json":
	}

	// case ext is .tools => readzip, pass to json
	// case ext is .json => read
	// case else => unsupported error

	return Library{}, errors.New("unsupported library file")
}

func ImportLibrary(name, path string) error {
	l := config.Logger.WithFields("fn", "tool.ImportLibrary")

	if err := CleanLibraryCache(); err != nil {
		return err
	}
	_, err := LoadLibrary(path)
	if err != nil {
		return err
	}

	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	targetName := ksuid.New().String() + filepath.Ext(path)
	targetPath := filepath.Join(config.LibraryCachePath(), targetName)

	if err := os.MkdirAll(config.LibraryCachePath(), 0750); err != nil {
		return err
	}

	dst, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	l.Debug("copying tool library to cache", "src", path, "dst", targetPath)
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	l.Debug("adding library to config", "name", name)
	if err := config.AddLibrary(name, targetPath); err != nil {

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
func CleanLibraryCache() error {
	l := config.Logger.WithFields("fn", "tool.CleanLibraryCache")

	l.Debug("loading stored libraries")
	libs, err := config.Libraries()
	if err != nil {
		return err
	}

	validPaths := make(map[string]string)

	for _, lib := range libs {
		validPaths[filepath.Base(lib.Path)] = lib.Name

		if _, err := os.Stat(lib.Path); os.IsNotExist(err) {
			l.Warn("removing invalid library from the config", "name", lib.Name, "path", lib.Path)
			if err := config.RemoveLibrary(lib.Name); err != nil {
				l.WithError(err).Error("cannot remove invalid library", "name", lib.Name)
			}

		}
	}

	l.Debug("Reading library cache", "path", config.LibraryCachePath())
	files, err := ioutil.ReadDir(config.LibraryCachePath())
	if err != nil {
		return err
	}

	for _, f := range files {
		if _, ok := validPaths[f.Name()]; !ok {
			path := filepath.Join(config.LibraryCachePath(), f.Name())
			l.Warn("removing orphan library", "path", path)
			os.Remove(path)
		}
	}

	return nil
}
