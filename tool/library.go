package tool

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
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

func LoadLibrary(path string) (Library, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return Library{}, err
	}
	if stat.IsDir() {
		return Library{}, errors.New("Library is directory")
	}

	// if not file, return error
	// if not exists, return error

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

/*

tool.LoadLibrary(somepath) => Library
tool.LoadLibrary(anotherPath) => Library

*/
