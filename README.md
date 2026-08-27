# commitkit

Go-native commit message parser and linter.

**Conventional Commits is a policy, not part of the parser.**

[Русская версия](README.ru.md)

## Architecture

Dependency direction is inward — core does not depend on adapters.

| Layer       | Packages                                                                | Role                              |
| ----------- | ----------------------------------------------------------------------- | --------------------------------- |
| Core        | `internal/commit`, `internal/lint`, `internal/policy`, `internal/lexer` | Domain logic                      |
| Application | `cli`, `internal/config`, `internal/hook`                               | CLI, config, hooks                |
| Adapters    | `commit`, `lint`, `policy`                                              | Public API for external importers |
| Entrypoint  | `main.go`                                                               | `main` only                       |

- **Core never imports** public packages or `cli`.
- **Public packages** are thin re-exports (type aliases + forwarding functions).
- **`cli` lives outside `internal`** — it is an application adapter, not library core.

## Features

- Public domain adapters + internal core
- Lexer + parser with source positions
- Git trailer support (continuations, multiline, BREAKING CHANGE)
- Composable rules / rule engine
- `.commitkit.yml` configuration (`gopkg.in/yaml.v3`)
- Git `commit-msg` hook installer
- Suitable for hooks and CI

## Configuration

`.commitkit.yml` (searched upward from cwd):

```yaml
types:
  - feat
  - fix
  - docs
  - chore

description:
  min: 1
  max: 72
  lowercase: true

scope:
  required: false
  enum: [api, ui]

header:
  max_length: 100

rules:
  no_trailing_period: true
  breaking_change_footer: true
```

No file -> built-in Conventional Commits defaults.

```bash
commitkit check --config path/to/.commitkit.yml
```

## CLI

```bash
go install github.com/destyk/commitkit@latest

# or build and install
make setup

git log -1 --format=%B | commitkit check
commitkit check --file .git/COMMIT_EDITMSG
commitkit parse
commitkit install-hook
commitkit uninstall-hook
```

| Exit | Meaning                     |
| ---- | --------------------------- |
| 0    | Success                     |
| 1    | Parse / policy / hook error |
| 2    | Usage / I/O / config error  |

## Library usage

```go
import (
    "github.com/destyk/commitkit/commit"
    "github.com/destyk/commitkit/lint"
    "github.com/destyk/commitkit/policy"
)

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

## Custom policy

Compose built-in rules or implement your own.

### From built-in rules

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

### Custom rule

Any type that implements `lint.Rule` can be used:

```go
package main

import (
    "strings"

    "github.com/destyk/commitkit/commit"
    "github.com/destyk/commitkit/lint"
    "github.com/destyk/commitkit/policy"
)

// ticketRefRule requires a footer like "Refs: PROJ-123".
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

func main() {
    message, err := commit.Parse(`feat(api): add pagination

Refs: API-42`)
    if err != nil {
        panic(err)
    }

    rules := policy.Custom(
        policy.TypeEnum("feat", "fix"),
        policy.RequireScope(),
        ticketRefRule{},
    )

    result := lint.Check(message, rules)
    if !result.Valid() {
        for _, v := range result.Violations {
            println(v.Rule + ": " + v.Message)
        }
    }
}
```

Via config, project-specific constraints are usually enough (types, scope enum, lengths). Code-level custom rules are for checks that do not map cleanly to YAML.

## Development

```bash
make test
make lint
gofmt -w .
make build
```

## License

MIT
