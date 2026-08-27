package cli

import (
	"bytes"
	"strings"
	"testing"
)

var testBuildInfo = BuildInfo{
	Version:   "test",
	Commit:    "test",
	BuildDate: "test",
}

func TestRunCheckValid(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := Run(
		[]string{"check"},
		strings.NewReader("feat(api): add pagination\n"),
		&out,
		&errOut,
		testBuildInfo,
	)

	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, errOut.String())
	}
	if !strings.Contains(out.String(), "commit message is valid") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestRunCheckInvalid(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := Run(
		[]string{"check"},
		strings.NewReader("feat: Add pagination\n"),
		&out,
		&errOut,
		testBuildInfo,
	)

	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if !strings.Contains(out.String(), "description-lowercase") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestRunHelp(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := Run(
		[]string{"help"},
		strings.NewReader(""),
		&out,
		&errOut,
		testBuildInfo,
	)

	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	if !strings.Contains(out.String(), "Usage:") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestRunParse(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := Run(
		[]string{"parse"},
		strings.NewReader("feat(api)!: break things\n\nBREAKING CHANGE: yes\n"),
		&out,
		&errOut,
		testBuildInfo,
	)

	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, errOut.String())
	}
	if !strings.Contains(out.String(), "Type: feat") {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "Breaking: true") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
