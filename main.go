package main

import (
	"os"

	"github.com/destyk/commitkit/cli"
)

// It is injected at build time via -ldflags.
var (
	// Version is the application version.
	version = "dev"

	// Commit is the Git commit used to build the binary.
	commit = "unknown"

	// BuildDate is the UTC build time.
	buildDate = "unknown"
)

func main() {
	r := cli.Run(
		os.Args[1:],
		os.Stdin,
		os.Stdout,
		os.Stderr,
		cli.BuildInfo{
			Version:   version,
			Commit:    commit,
			BuildDate: buildDate,
		},
	)
	os.Exit(r)
}
