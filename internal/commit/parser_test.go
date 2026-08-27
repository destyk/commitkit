package commit

import (
	"strings"
	"testing"
)

func TestParseSimpleCommit(t *testing.T) {
	message, err := Parse("feat: add pagination")
	if err != nil {
		t.Fatal(err)
	}

	if message.Header.Type != "feat" {
		t.Fatalf("type = %q, want %q", message.Header.Type, "feat")
	}
	if message.Header.Scope != "" {
		t.Fatalf("scope = %q, want empty", message.Header.Scope)
	}
	if message.Header.Breaking {
		t.Fatal("breaking = true, want false")
	}
	if message.Header.Description != "add pagination" {
		t.Fatalf("description = %q, want %q", message.Header.Description, "add pagination")
	}
	if !message.Header.TypeSpan.IsValid() {
		t.Fatal("TypeSpan should be valid")
	}
	if !message.Header.DescSpan.IsValid() {
		t.Fatal("DescSpan should be valid")
	}
}

func TestParseScopedCommit(t *testing.T) {
	message, err := Parse("feat(api): add pagination")
	if err != nil {
		t.Fatal(err)
	}
	if message.Header.Type != "feat" {
		t.Fatalf("type = %q", message.Header.Type)
	}
	if message.Header.Scope != "api" {
		t.Fatalf("scope = %q", message.Header.Scope)
	}
	if !message.Header.ScopeSpan.IsValid() {
		t.Fatal("ScopeSpan should be valid")
	}
}

func TestParseBreakingCommit(t *testing.T) {
	message, err := Parse("feat(api)!: remove old endpoint")
	if err != nil {
		t.Fatal(err)
	}
	if !message.Header.Breaking {
		t.Fatal("breaking = false, want true")
	}
	if !message.Header.BreakingSpan.IsValid() {
		t.Fatal("BreakingSpan should be valid")
	}
}

func TestParseBody(t *testing.T) {
	message, err := Parse(`feat(api): add pagination

Add cursor based pagination.

This allows clients to paginate large collections.`)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Add cursor based pagination.\n\nThis allows clients to paginate large collections."
	if message.Body.Text != expected {
		t.Fatalf("body = %q, want %q", message.Body.Text, expected)
	}
}

func TestParseFooters(t *testing.T) {
	message, err := Parse(`feat(api)!: remove endpoint

Remove the old endpoint.

BREAKING CHANGE: endpoint was removed.
Refs: MOV-123`)
	if err != nil {
		t.Fatal(err)
	}

	if len(message.Footers) != 2 {
		t.Fatalf("footers = %d, want 2", len(message.Footers))
	}
	if message.Footers[0].Token != "BREAKING CHANGE" {
		t.Fatalf("footer token = %q", message.Footers[0].Token)
	}
	if message.Footers[0].Value != "endpoint was removed." {
		t.Fatalf("footer value = %q", message.Footers[0].Value)
	}
	if message.Footers[1].Token != "Refs" {
		t.Fatalf("footer token = %q", message.Footers[1].Token)
	}
	if !message.HasBreakingFooter() {
		t.Fatal("expected HasBreakingFooter() == true")
	}
}

func TestParseFooterContinuation(t *testing.T) {
	message, err := Parse(`fix: handle edge case

This is the body.

Reviewed-by: Alice <alice@example.com>
    Bob <bob@example.com>
Refs: #42`)
	if err != nil {
		t.Fatal(err)
	}

	if len(message.Footers) != 2 {
		t.Fatalf("footers = %d, want 2; got %+v", len(message.Footers), message.Footers)
	}

	if message.Footers[0].Token != "Reviewed-by" {
		t.Fatalf("token = %q", message.Footers[0].Token)
	}
	if !strings.Contains(message.Footers[0].Value, "Bob <bob@example.com>") {
		t.Fatalf("continuation missing: %q", message.Footers[0].Value)
	}
	if message.Footers[1].Token != "Refs" {
		t.Fatalf("token = %q", message.Footers[1].Token)
	}
}

func TestParseFooterMultilineValue(t *testing.T) {
	message, err := Parse(`feat: something

BREAKING CHANGE: the old API was removed.
    Clients must migrate to the new endpoint
    before upgrading.
Refs: ABC-1`)
	if err != nil {
		t.Fatal(err)
	}

	if len(message.Footers) != 2 {
		t.Fatalf("footers = %d, want 2", len(message.Footers))
	}
	val := message.Footers[0].Value
	if !strings.Contains(val, "Clients must migrate") {
		t.Fatalf("multiline value incomplete: %q", val)
	}
	if !strings.Contains(val, "before upgrading.") {
		t.Fatalf("multiline value incomplete: %q", val)
	}
}

func TestBodyWithColonNotTreatedAsTrailer(t *testing.T) {
	message, err := Parse(`feat: improve docs

Use foo: bar when needed.

This is still body.`)
	if err != nil {
		t.Fatal(err)
	}

	if len(message.Footers) != 0 {
		t.Fatalf("expected no footers, got %+v", message.Footers)
	}
	if !strings.Contains(message.Body.Text, "Use foo: bar") {
		t.Fatalf("body lost the colon line: %q", message.Body.Text)
	}
}

func TestParseCRLF(t *testing.T) {
	message, err := Parse("feat: add thing\r\n\r\nBody\r\n")
	if err != nil {
		t.Fatal(err)
	}
	if message.Body.Text != "Body" {
		t.Fatalf("body = %q, want %q", message.Body.Text, "Body")
	}
}

func TestParseRejectsEmptyMessage(t *testing.T) {
	_, err := Parse("   \n\n")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseRejectsMissingColon(t *testing.T) {
	_, err := Parse("feat add pagination")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "header must contain ':'") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsEmptyDescription(t *testing.T) {
	_, err := Parse("feat:")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseRejectsEmptyScope(t *testing.T) {
	_, err := Parse("feat(): add pagination")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseRejectsMalformedScope(t *testing.T) {
	_, err := Parse("feat(api: add pagination")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseHeaderPositions(t *testing.T) {
	msg, err := Parse("feat(api)!: add pagination")
	if err != nil {
		t.Fatal(err)
	}

	if msg.Header.TypeSpan.Start.Column != 1 {
		t.Fatalf("TypeSpan start column = %d", msg.Header.TypeSpan.Start.Column)
	}
	if msg.Header.BreakingSpan.Start.Column == 0 {
		t.Fatal("BreakingSpan not set")
	}
}
