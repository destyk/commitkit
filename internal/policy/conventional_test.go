package policy

import (
	"testing"

	"github.com/destyk/commitkit/internal/commit"
	"github.com/destyk/commitkit/internal/lint"
)

func TestConventionalCommitsAcceptsValidCommit(t *testing.T) {
	message, err := commit.Parse("feat(api): add pagination")
	if err != nil {
		t.Fatal(err)
	}

	result := lint.Check(message, ConventionalCommits())
	if !result.Valid() {
		t.Fatalf("unexpected violations: %+v", result.Violations)
	}
}

func TestConventionalCommitsRejectsUnknownType(t *testing.T) {
	message, err := commit.Parse("feature(api): add pagination")
	if err != nil {
		t.Fatal(err)
	}

	result := lint.Check(message, ConventionalCommits())
	if result.Valid() {
		t.Fatal("expected violation")
	}

	if result.Violations[0].Rule != "type-enum" {
		t.Fatalf("rule = %q, want %q", result.Violations[0].Rule, "type-enum")
	}
}

func TestConventionalCommitsRejectsUppercaseDescription(t *testing.T) {
	message, err := commit.Parse("feat: Add pagination")
	if err != nil {
		t.Fatal(err)
	}

	result := lint.Check(message, ConventionalCommits())
	if result.Valid() {
		t.Fatal("expected violation")
	}
}

func TestConventionalCommitsRejectsLongDescription(t *testing.T) {
	message, err := commit.Parse(
		"feat: " +
			"this is an intentionally very long commit description " +
			"that exceeds the configured seventy two character limit",
	)
	if err != nil {
		t.Fatal(err)
	}

	result := lint.Check(message, ConventionalCommits())
	if result.Valid() {
		t.Fatal("expected violation")
	}
}

func TestConventionalCommitsRequiresBreakingFooter(t *testing.T) {
	message, err := commit.Parse("feat(api)!: remove endpoint")
	if err != nil {
		t.Fatal(err)
	}

	result := lint.Check(message, ConventionalCommits())
	if result.Valid() {
		t.Fatal("expected violation")
	}

	found := false
	for _, v := range result.Violations {
		if v.Rule == "breaking-change-footer" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected breaking-change-footer violation")
	}
}

func TestConventionalCommitsAcceptsBreakingFooter(t *testing.T) {
	message, err := commit.Parse(`feat(api)!: remove endpoint

BREAKING CHANGE: endpoint was removed.`)
	if err != nil {
		t.Fatal(err)
	}

	result := lint.Check(message, ConventionalCommits())
	if !result.Valid() {
		t.Fatalf("unexpected violations: %+v", result.Violations)
	}
}

func TestConventionalCommitsAcceptsBreakingHyphenFooter(t *testing.T) {
	message, err := commit.Parse(`feat(api)!: remove endpoint

BREAKING-CHANGE: endpoint was removed.`)
	if err != nil {
		t.Fatal(err)
	}

	result := lint.Check(message, ConventionalCommits())
	if !result.Valid() {
		t.Fatalf("unexpected violations: %+v", result.Violations)
	}
}

func TestEngine(t *testing.T) {
	message, err := commit.Parse("feat: add thing")
	if err != nil {
		t.Fatal(err)
	}

	engine := lint.NewEngine(ConventionalCommits()...)
	result := engine.Check(message)
	if !result.Valid() {
		t.Fatalf("unexpected: %+v", result.Violations)
	}
}
