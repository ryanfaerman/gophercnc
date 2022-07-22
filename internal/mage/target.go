package mage

import (
	"fmt"
	"path/filepath"
)

type Target struct {
	GOOS      string
	GOARCH    string
	SourceDir string
}

func (t Target) Name() string {
	n := filepath.Base(ModulePath())
	if t.SourceDir != "" {
		n = t.SourceDir
	}
	name := fmt.Sprintf(
		"%s-%s-%s",
		n,
		t.GOOS,
		t.GOARCH,
	)

	if t.GOOS == "windows" {
		name += ".exe"
	}

	return name
}
