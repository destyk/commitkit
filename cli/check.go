package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/destyk/commitkit/internal/commit"
	"github.com/destyk/commitkit/internal/config"
	"github.com/destyk/commitkit/internal/lint"
)

func check(
	args []string,
	in io.Reader,
	out io.Writer,
	errOut io.Writer,
) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(errOut)

	file := fs.String("file", "", "commit message file")
	configPath := fs.String("config", "", "path to .commitkit.yml (default: search upwards from cwd)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	data, err := readInput(*file, in)
	if err != nil {
		fmt.Fprintf(errOut, "read commit message: %v\n", err)
		return 2
	}

	message, err := commit.Parse(string(data))
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}

	var cfgResult config.LoadResult
	if *configPath != "" {
		cfgResult, err = config.Load(*configPath)
	} else {
		cfgResult, err = config.FindAndLoad("")
	}
	if err != nil {
		fmt.Fprintf(errOut, "load config: %v\n", err)
		return 2
	}

	result := lint.Check(message, cfgResult.Config.ToRules())

	for _, violation := range result.Violations {
		printViolation(out, violation)
	}

	if result.Valid() {
		fmt.Fprintln(out, "commit message is valid")
		return 0
	}

	return 1
}

func printViolation(out io.Writer, v lint.Violation) {
	loc := v.Position()
	if loc != "" {
		fmt.Fprintf(out, "%s:%s: %s\n", loc, v.Rule, v.Message)
	} else {
		fmt.Fprintf(out, "%s: %s\n", v.Rule, v.Message)
	}
}
