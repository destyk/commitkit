package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/destyk/commitkit/internal/hook"
)

func installHook(
	args []string,
	in io.Reader,
	out io.Writer,
	errOut io.Writer,
) int {
	fs := flag.NewFlagSet("install-hook", flag.ContinueOnError)
	fs.SetOutput(errOut)

	force := fs.Bool("force", false, "overwrite an existing non-commitkit hook")
	path := fs.String("path", "", "path to the git repository root (default: discover from cwd)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	res, err := hook.Install(hook.Options{
		RepoRoot: *path,
		Force:    *force,
	})
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}

	fmt.Fprintf(out, "%s\n  %s\n", res.Message, res.HookPath)
	return 0
}

func uninstallHook(
	args []string,
	in io.Reader,
	out io.Writer,
	errOut io.Writer,
) int {
	fs := flag.NewFlagSet("uninstall-hook", flag.ContinueOnError)
	fs.SetOutput(errOut)

	force := fs.Bool("force", false, "remove even if the hook is not managed by commitkit")
	path := fs.String("path", "", "path to the git repository root (default: discover from cwd)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	res, err := hook.Uninstall(hook.Options{
		RepoRoot: *path,
		Force:    *force,
	})
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}

	fmt.Fprintf(out, "%s\n  %s\n", res.Message, res.HookPath)
	return 0
}
