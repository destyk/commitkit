package cli

import (
	"fmt"
	"io"
	"os"
)

type BuildInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

// Run executes the commitkit CLI.
//
// The returned integer is intended to be passed to os.Exit.
func Run(
	args []string,
	in io.Reader,
	out io.Writer,
	errOut io.Writer,
	buildInfo BuildInfo,
) int {
	if len(args) == 0 {
		usage(out)
		return 2
	}

	switch args[0] {
	case "check":
		return check(args[1:], in, out, errOut)

	case "parse":
		return parse(args[1:], in, out, errOut)

	case "install-hook":
		return installHook(args[1:], in, out, errOut)

	case "uninstall-hook":
		return uninstallHook(args[1:], in, out, errOut)

	case "version":
		fmt.Fprintf(
			out,
			"commitkit %s (commit: %s built: %s)",
			buildInfo.Version,
			buildInfo.Commit,
			buildInfo.BuildDate,
		)
		return 0

	case "help", "-h", "--help":
		usage(out)
		return 0

	default:
		fmt.Fprintf(errOut, "unknown command %q\n\n", args[0])
		usage(errOut)
		return 2
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, `commitkit - commit message parser and linter

Usage:
  commitkit check [--file FILE] [--config FILE]
  commitkit parse [--file FILE]
  commitkit install-hook [--force] [--path DIR]
  commitkit uninstall-hook [--force] [--path DIR]
  commitkit version

Commands:
  check            Validate a commit message.
  parse            Parse and print a structured commit message.
  install-hook     Install the git commit-msg hook.
  uninstall-hook   Remove the git commit-msg hook.
  version          Print the version.

Configuration:
  check loads .commitkit.yml by walking up from the current directory.
  Use --config to pass an explicit path. If no config is found, the
  built-in Conventional Commits policy is used.

Input:
  If --file is omitted, stdin is used.

Examples:
  git log -1 --format=%B | commitkit check
  commitkit check --file .git/COMMIT_EDITMSG
  commitkit check --config .commitkit.yml
  commitkit install-hook`)
}

func readInput(file string, in io.Reader) ([]byte, error) {
	if file != "" {
		return os.ReadFile(file)
	}
	return io.ReadAll(in)
}
