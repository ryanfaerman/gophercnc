//go:build mage
// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	helpers "github.com/ryanfaerman/gophercnc/internal/mage"
)

var (
	goexe = "go"
	dirs  = []string{"bin", "tmp"}

	targets = []helpers.Target{
		{GOOS: "linux", GOARCH: "amd64"},
		{GOOS: "darwin", GOARCH: "amd64"},
		{GOOS: "windows", GOARCH: "amd64"},
	}
)

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}
}

func flags() string {
	f := helpers.LDFlags{}

	f["github.com/ryanfaerman/gophercnc/version.ApplicationName"] = filepath.Base(helpers.ModulePath())
	f["github.com/ryanfaerman/gophercnc/version.BuildDate"] = time.Now().Format(time.RFC3339)
	f["github.com/ryanfaerman/gophercnc/version.BuildTag"] = helpers.GitTag()
	f["github.com/ryanfaerman/gophercnc/version.CommitHash"] = helpers.GitCommitHash()

	return f.String()
}

func ensureDirs() error {
	fmt.Println("--> Ensuring output directories")

	for _, dir := range dirs {
		if !helpers.FileExists("./" + dir) {
			fmt.Printf("    creating './%s'\n", dir)
			if err := os.MkdirAll("./"+dir, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

// Clean up after yourself
func Clean() {
	fmt.Println("--> Cleaning output directories")

	for _, dir := range dirs {
		fmt.Printf("    removing './%s'\n", dir)
		os.RemoveAll("./" + dir)
	}
}

// Vendor dependencies with go modules
func Vendor() {
	fmt.Println("--> Updating dependencies")
	sh.Run(goexe, "mod", "tidy")
}

func commands() []string {
	c := []string{}

	if files, err := ioutil.ReadDir("./cmd"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				c = append(c, file.Name())
			}
		}
	}

	return c
}

// Build the application for local running
func Build() error {
	mg.SerialDeps(Vendor, ensureDirs)

	for _, command := range commands() {
		fmt.Printf("--> Building '%s'\n", command)

		binaryPath := filepath.Join("./bin", command)
		sourcePath := filepath.Join(helpers.ModulePath(), "/cmd", command)
		if err := sh.Run(goexe, "build", "-o", binaryPath, "-ldflags="+flags(), "-trimpath", sourcePath); err != nil {
			return err
		}
	}

	return nil

}

// Release the application for all defined targets
func Release() error {
	mg.SerialDeps(Vendor, ensureDirs, Test)

	cmds := commands()

	var wg sync.WaitGroup
	wg.Add(len(targets) * len(cmds))
	for _, c := range cmds {
		fmt.Printf("--> Building '%s' for release\n", c)
		for _, t := range targets {
			t.SourceDir = c
			go func(t helpers.Target) {
				defer wg.Done()

				fmt.Printf("      Building %s\n", t.Name())

				env := map[string]string{
					"GOOS":   t.GOOS,
					"GOARCH": t.GOARCH,
				}

				binaryPath := filepath.Join("./bin", t.Name())
				sourcePath := filepath.Join(helpers.ModulePath(), "/cmd", t.SourceDir)

				err := sh.RunWith(env, goexe, "build", "-o", binaryPath, "-ldflags="+flags(), "-trimpath", sourcePath)
				if err != nil {
					fmt.Printf("compilation failed: %s\n", err.Error())
					return
				}

			}(t)
		}
	}
	wg.Wait()

	return nil
}

// Lint the codebase, checking for common errors
func Lint() {
	fmt.Println("--> Linting codebase")

	c := exec.Command("gometalinter", "-e", "internal", "-e", "go/pkg/mod", "./...")
	c.Env = os.Environ()
	out, err := c.CombinedOutput()
	if err == nil {
		fmt.Println("    no issues detected")
	} else {
		fmt.Print("    ")
		fmt.Println(strings.Replace(string(out), "\n", "\n    ", -1))
	}
}

// Test the codebase
func Test() error {
	mg.SerialDeps(Vendor, ensureDirs)

	fmt.Println("--> Testing codebase")
	results, err := sh.Output(goexe, "test", "-cover", "./...")
	fmt.Print("    ")
	fmt.Println(strings.Replace(results, "\n", "\n    ", -1))

	return err
}
