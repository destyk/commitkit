package hook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallAndUninstall(t *testing.T) {
	tmp := t.TempDir()
	// Fake a git repo.
	if err := os.MkdirAll(filepath.Join(tmp, ".git", "hooks"), 0o755); err != nil {
		t.Fatal(err)
	}

	res, err := Install(Options{RepoRoot: tmp})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Created {
		t.Fatalf("expected Created, got %+v", res)
	}

	data, err := os.ReadFile(res.HookPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Managed by commitkit") {
		t.Fatalf("hook content missing marker: %s", data)
	}

	// Executable?
	info, err := os.Stat(res.HookPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&0o111 == 0 {
		t.Fatal("hook is not executable")
	}

	// Second install should skip.
	res2, err := Install(Options{RepoRoot: tmp})
	if err != nil {
		t.Fatal(err)
	}
	if !res2.Skipped {
		t.Fatalf("expected Skipped, got %+v", res2)
	}

	// Uninstall.
	ures, err := Uninstall(Options{RepoRoot: tmp})
	if err != nil {
		t.Fatal(err)
	}
	if ures.Skipped {
		t.Fatal("expected removal")
	}
	if _, err := os.Stat(res.HookPath); !os.IsNotExist(err) {
		t.Fatal("hook file still exists")
	}
}

func TestInstallRefusesForeignHook(t *testing.T) {
	tmp := t.TempDir()

	hooksDir := filepath.Join(tmp, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatal(err)
	}

	hookPath := filepath.Join(hooksDir, "commit-msg")
	if err := os.WriteFile(hookPath, []byte("#!/bin/sh\necho foreign\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := Install(Options{RepoRoot: tmp})
	if err == nil {
		t.Fatal("expected error for foreign hook")
	}

	// Force should overwrite.
	res, err := Install(Options{RepoRoot: tmp, Force: true})
	if err != nil {
		t.Fatal(err)
	}

	if !res.Updated {
		t.Fatalf("expected Updated, got %+v", res)
	}
}
