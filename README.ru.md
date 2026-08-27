# commitkit

Парсер и линтер commit-сообщений на чистом Go.

**Conventional Commits — это политика, а не часть парсера.**

[English version](README.md)

## Архитектура

Зависимости направлены внутрь: ядро не зависит от адаптеров.

| Слой        | Пакеты                                                                  | Роль                               |
| ----------- | ----------------------------------------------------------------------- | ---------------------------------- |
| Ядро        | `internal/commit`, `internal/lint`, `internal/policy`, `internal/lexer` | Доменная логика                    |
| Приложение  | `cli`, `internal/config`, `internal/hook`                               | CLI, конфиг, хуки                  |
| Адаптеры    | `commit`, `lint`, `policy`                                              | Публичный API для внешних импортов |
| Точка входа | `main.go`                                                               | Только `main`                      |

- **Ядро никогда не импортирует** публичные пакеты или `cli`.
- **Публичные пакеты** — тонкие реэкспорты (type aliases и прокси-функции).
- **`cli` живёт вне `internal`** — это прикладной адаптер, а не библиотечное ядро.

## Возможности

- Публичные доменные адаптеры и внутреннее ядро
- Лексер и парсер с позициями в исходнике (spans)
- Поддержка Git trailers (continuation, multiline, BREAKING CHANGE)
- Компонуемые правила и rule engine
- Конфигурация `.commitkit.yml` (`gopkg.in/yaml.v3`)
- Установщик git-хука `commit-msg`
- Удобно для хуков и CI

## Конфигурация

`.commitkit.yml` (поиск вверх от текущего каталога):

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

Если файла нет — встроенные значения Conventional Commits.

```bash
commitkit check --config path/to/.commitkit.yml
```

## CLI

```bash
go install github.com/destyk/commitkit@latest

# или сборка и установка
make setup

git log -1 --format=%B | commitkit check
commitkit check --file .git/COMMIT_EDITMSG
commitkit parse
commitkit install-hook
commitkit uninstall-hook
```

| Код | Значение                         |
| --- | -------------------------------- |
| 0   | Успех                            |
| 1   | Ошибка разбора / политики / хука |
| 2   | Ошибка usage / I/O / конфига     |

## Использование как библиотеки

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

## Кастомная политика

Можно собирать правила из встроенных или писать свои.

### Из встроенных правил

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

### Своё правило

Достаточно реализовать интерфейс `lint.Rule`:

```go
package main

import (
    "strings"

    "github.com/destyk/commitkit/commit"
    "github.com/destyk/commitkit/lint"
    "github.com/destyk/commitkit/policy"
)

// ticketRefRule требует footer вида "Refs: PROJ-123".
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

Через конфиг обычно хватает ограничений проекта (types, scope enum, длины). Правила в коде нужны для проверок, которые плохо выражаются в YAML.

## Разработка

```bash
make test
make lint
gofmt -w .
make build
```

## Лицензия

MIT
