package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/destyk/commitkit/internal/commit"
	"github.com/destyk/commitkit/internal/lint"
)

func TestLoadAndRules(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, FileName)

	content := `
types:
  - feat
  - fix

description:
  min: 3
  max: 50
  lowercase: true

scope:
  required: true
  enum:
    - api
    - ui

header:
  max_length: 80

rules:
  no_trailing_period: true
  breaking_change_footer: true
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Found {
		t.Fatal("expected Found")
	}

	rules := res.Config.ToRules()

	// Valid commit.
	msg, err := commit.Parse("feat(api): add pagination")
	if err != nil {
		t.Fatal(err)
	}
	result := lint.Check(msg, rules)
	if !result.Valid() {
		t.Fatalf("expected valid, got %+v", result.Violations)
	}

	// Unknown type.
	msg, err = commit.Parse("chore(api): touch files")
	if err != nil {
		t.Fatal(err)
	}
	result = lint.Check(msg, rules)
	if result.Valid() {
		t.Fatal("expected type-enum violation")
	}

	// Missing scope.
	msg, err = commit.Parse("feat: add pagination")
	if err != nil {
		t.Fatal(err)
	}
	result = lint.Check(msg, rules)
	if result.Valid() {
		t.Fatal("expected scope-required violation")
	}

	// Bad scope.
	msg, err = commit.Parse("feat(db): add pagination")
	if err != nil {
		t.Fatal(err)
	}

	result = lint.Check(msg, rules)
	if result.Valid() {
		t.Fatal("expected scope-enum violation")
	}
}

func TestFindAndLoadWalksUp(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(root, FileName)
	if err := os.WriteFile(cfgPath, []byte("types:\n  - feat\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := FindAndLoad(nested)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Found {
		t.Fatal("expected to find config by walking up")
	}
	if len(res.Config.Types) != 1 || res.Config.Types[0] != "feat" {
		t.Fatalf("types = %v", res.Config.Types)
	}
}

func TestFindAndLoadDefaults(t *testing.T) {
	dir := t.TempDir()
	res, err := FindAndLoad(dir)
	if err != nil {
		t.Fatal(err)
	}
	if res.Found {
		t.Fatal("expected no config")
	}

	// Defaults should accept a normal conventional commit.
	msg, err := commit.Parse("feat: add thing")
	if err != nil {
		t.Fatal(err)
	}
	result := lint.Check(msg, res.Config.ToRules())
	if !result.Valid() {
		t.Fatalf("defaults rejected valid commit: %+v", result.Violations)
	}
}

func TestDisableLowercase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, FileName)
	content := `
description:
  lowercase: false
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := commit.Parse("feat: Add Thing")
	if err != nil {
		t.Fatal(err)
	}

	result := lint.Check(msg, res.Config.ToRules())
	if !result.Valid() {
		t.Fatalf("expected valid when lowercase disabled: %+v", result.Violations)
	}
}
