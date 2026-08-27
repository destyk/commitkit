package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/destyk/commitkit/internal/commit"
)

func parse(
	args []string,
	in io.Reader,
	out io.Writer,
	errOut io.Writer,
) int {
	fs := flag.NewFlagSet("parse", flag.ContinueOnError)
	fs.SetOutput(errOut)

	file := fs.String("file", "", "commit message file")

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

	fmt.Fprintf(out, "Type: %s\n", message.Header.Type)
	fmt.Fprintf(out, "Scope: %s\n", message.Header.Scope)
	fmt.Fprintf(out, "Breaking: %t\n", message.Header.Breaking)
	fmt.Fprintf(out, "Description: %s\n", message.Header.Description)

	if message.Header.TypeSpan.IsValid() {
		fmt.Fprintf(out, "TypeSpan: %s\n", message.Header.TypeSpan)
	}
	if message.Header.ScopeSpan.IsValid() {
		fmt.Fprintf(out, "ScopeSpan: %s\n", message.Header.ScopeSpan)
	}
	if message.Header.BreakingSpan.IsValid() {
		fmt.Fprintf(out, "BreakingSpan: %s\n", message.Header.BreakingSpan)
	}
	if message.Header.DescSpan.IsValid() {
		fmt.Fprintf(out, "DescSpan: %s\n", message.Header.DescSpan)
	}

	if message.Body.Text != "" {
		fmt.Fprintf(out, "\nBody:\n%s\n", message.Body.Text)
	}

	if len(message.Footers) > 0 {
		fmt.Fprintln(out, "\nFooters:")
		for _, footer := range message.Footers {
			fmt.Fprintf(out, "  %s: %s\n", footer.Token, footer.Value)
		}
	}

	return 0
}
