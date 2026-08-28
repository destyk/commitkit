# CONTRIBUTING

This document is for contributors and anyone embedding commitkit as a library.
For a product overview and “why use this”, see [README.md](README.md).

[Русская версия](CONTRIBUTING.ru.md)

---

## Architecture

Dependency direction is inward — `internal` does not depend on adapters.

| Layer       | Packages                                                                | Role                              |
| ----------- | ----------------------------------------------------------------------- | --------------------------------- |
| Core        | `internal/commit`, `internal/lint`, `internal/policy`, `internal/lexer` | Domain logic                      |
| Application | `cli`, `internal/config`, `internal/hook`                               | CLI, config, hooks                |
| Adapters    | `commit`, `lint`, `policy`                                              | Public API for external importers |
| Entrypoint  | `main.go`                                                               | `main` only                       |

- **Core never imports** public packages or `cli`.
- **Public packages** are thin re-exports (type aliases + forwarding functions).
- **`cli` lives outside `internal`** — it is an application adapter, not library core.
- **Config is internal only** — the CLI loads `.commitkit.yml`; library users
  compose `policy` / `lint` rules in code.

---

## Development setup

```bash
git clone https://github.com/destyk/commitkit.git
cd commitkit

make setup
make test
make lint
gofmt -w .
make build
```

Pull requests should keep the dependency rule: `internal` no adapters, no `cli`.

---

## Library usage

Import only the public domain packages:

```go
import (
    "github.com/destyk/commitkit/commit"
    "github.com/destyk/commitkit/lint"
    "github.com/destyk/commitkit/policy"
)
```

Do **not** import `internal/...`.

```go
message, err := commit.Parse("feat(api): add pagination")
if err != nil {
    panic(err)
}

result := lint.Check(message, policy.ConventionalCommits())
if !result.Valid() {
    for _, v := range result.Violations {
        fmt.Printf("%s: %s: %s\n", v.Position(), v.Rule, v.Message)
    }
}
```

---

## Built-in policy helpers

| Helper                               | Role                                      |
| ------------------------------------ | ----------------------------------------- |
| `policy.ConventionalCommits()`       | Default CC rule set                       |
| `policy.TypeEnum(...)`               | Allowed types                             |
| `policy.DescriptionLength(min, max)` | Description length in runes               |
| `policy.DescriptionLowercase()`      | First letter lowercase                    |
| `policy.RequireScope()`              | Non-empty scope                           |
| `policy.ScopeEnum(...)`              | Allowed scopes                            |
| `policy.HeaderMaxLength(n)`          | Header line max length                    |
| `policy.NoTrailingPeriod()`          | No trailing `.` in description            |
| `policy.BreakingChangeFooter()`      | `BREAKING CHANGE` footer when `!` is used |
| `policy.Custom(rules...)`            | Compose arbitrary rules                   |

Example:

```go
rules := policy.Custom(
    policy.TypeEnum("feat", "fix", "chore"),
    policy.DescriptionLength(1, 72),
    policy.DescriptionLowercase(),
    policy.RequireScope(),
    policy.ScopeEnum("api", "ui", "db"),
    policy.NoTrailingPeriod(),
    policy.HeaderMaxLength(100),
    policy.BreakingChangeFooter(),
)

result := lint.Check(message, rules)
```

---

## Adding a rule

### 1. Implement `lint.Rule`

```go
type Rule interface {
    Name() string
    Check(message commit.Message) []lint.Violation
}
```

Example:

```go
type ticketRefRule struct{}

func (ticketRefRule) Name() string { return "ticket-ref" }

func (ticketRefRule) Check(message commit.Message) []lint.Violation {
    for _, footer := range message.Footers {
        if strings.EqualFold(footer.Token, "Refs") && footer.Value != "" {
            return nil
        }
    }
    return []lint.Violation{{
        Rule:     "ticket-ref",
        Message:  "commit must include a Refs footer (e.g. Refs: PROJ-123)",
        Severity: lint.SeverityError,
    }}
}
```

### 2. Wire it in

- **Library users** — pass the rule into `policy.Custom(...)` or `lint.NewEngine(...)`.
- **Built-in policy** — add a constructor under `internal/policy` and re-export it
  from the public `policy` package (type alias / forwarding function only).
- **YAML config** — only if the rule maps cleanly to a config field; extend
  `internal/config` and `Config.ToRules()`. Prefer code-level rules when the
  check does not fit YAML.

### 3. Conventions

- Put domain logic in `internal/...`; public packages stay thin.
- Prefer precise `Violation` messages and positions when the parser provides them.
- Add tests next to the rule (`internal/policy/*_test.go` or `internal/lint`).

---

## Config loading (CLI)

```text
FindAndLoad(startDir)  ->  walk up for .commitkit.yml
Load(path)             ->  explicit file
Default()              ->  when no file found (Found=false)
Config.ToRules()       ->  lint.Rules with CC defaults for omitted fields
```

The CLI never depends on public `config` — there is none. Config is an
application detail under `internal/config`.

---

## CLI surface

| Command                           | Role                                       |
| --------------------------------- | ------------------------------------------ |
| `check`                           | Parse message + run policy rules           |
| `parse`                           | Print structured message (debug / tooling) |
| `install-hook` / `uninstall-hook` | Manage `commit-msg` hook                   |
| `version` / `help`                | Meta                                       |

Console formatting lives in `cli/`; `internal/*` returns data and errors.

---

## License

By contributing, you agree that your contributions are licensed under the MIT License.
