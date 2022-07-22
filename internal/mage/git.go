package mage

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func GitTag() string {
	var err error
	tagOrBranch, ok := os.LookupEnv("CI_COMMIT_REF_NAME")
	if !ok {
		tagOrBranch, err = Shellout("git", "describe", "--tags")
		if err != nil {
			ee, ok := errors.Cause(err).(*exec.ExitError)
			if ok && ee.Exited() {
				// probably no git tag
				return "dev"
			}
			return "dev"
		}
	}

	return strings.TrimSuffix(tagOrBranch, "\n")
}

func GitCommitHash() string {
	var err error

	hash, ok := os.LookupEnv("CI_COMMIT_SHA")
	if !ok {
		hash, err = Shellout("git", "rev-parse", "HEAD")
		if err != nil {
			return ""
		}
	}

	hash = strings.TrimSpace(hash)
	if len(hash) >= 8 {
		return strings.TrimSpace(hash)[0:8]
	}

	return ""
}

func GitBranch() string {
	branch, err := Shellout("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return ""
	}
	branch = strings.TrimSpace(branch)
	return branch
}
