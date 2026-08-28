# CONTRIBUTING

Документ для контрибьюторов и тех, кто встраивает commitkit как библиотеку.
Обзор продукта и "зачем это" — в [README.ru.md](README.ru.md).

[English version](CONTRIBUTING.md)

---

## Архитектура

Направление зависимостей внутрь - `internal` не зависит от адаптеров.

| Слой        | Пакеты                                                                  | Роль                               |
| ----------- | ----------------------------------------------------------------------- | ---------------------------------- |
| Core        | `internal/commit`, `internal/lint`, `internal/policy`, `internal/lexer` | Доменная логика                    |
| Application | `cli`, `internal/config`, `internal/hook`                               | CLI, конфиг, хуки                  |
| Adapters    | `commit`, `lint`, `policy`                                              | Публичный API для внешних импортов |
| Entrypoint  | `main.go`                                                               | только `main`                      |

- **Core никогда не импортирует** публичные пакеты или `cli`.
- **Публичные пакеты** — тонкие re-export (type aliases + forwarding).
- **`cli` снаружи `internal`** — application adapter, не библиотечное ядро.
- **Config только internal** — `.commitkit.yml` читает CLI; library-пользователи
  собирают `policy` / `lint` правила в коде.

---

## Окружение для разработки

```bash
git clone https://github.com/destyk/commitkit.git
cd commitkit

make setup
make test
make lint
gofmt -w .
make build
```

В PR соблюдайте правило зависимостей: `internal` без адаптеров и без `cli`.

---

## Использование как библиотеки

Импортируйте только публичные domain-пакеты:

```go
import (
    "github.com/destyk/commitkit/commit"
    "github.com/destyk/commitkit/lint"
    "github.com/destyk/commitkit/policy"
)
```

**Не** импортируйте `internal/...`.

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

## Встроенные policy-хелперы

| Хелпер                               | Роль                             |
| ------------------------------------ | -------------------------------- |
| `policy.ConventionalCommits()`       | Набор правил CC по умолчанию     |
| `policy.TypeEnum(...)`               | Допустимые types                 |
| `policy.DescriptionLength(min, max)` | Длина description в рунах        |
| `policy.DescriptionLowercase()`      | Первая буква в нижнем регистре   |
| `policy.RequireScope()`              | Непустой scope                   |
| `policy.ScopeEnum(...)`              | Допустимые scopes                |
| `policy.HeaderMaxLength(n)`          | Макс. длина строки header        |
| `policy.NoTrailingPeriod()`          | Без `.` в конце description      |
| `policy.BreakingChangeFooter()`      | Footer `BREAKING CHANGE` при `!` |
| `policy.Custom(rules...)`            | Сборка произвольных правил       |

Пример:

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

## Добавление правила

### 1. Реализовать `lint.Rule`

```go
type Rule interface {
    Name() string
    Check(message commit.Message) []lint.Violation
}
```

Пример:

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

### 2. Подключить

- **Library** — передать правило в `policy.Custom(...)` или `lint.NewEngine(...)`.
- **Встроенная policy** — конструктор в `internal/policy` + thin re-export в
  публичном `policy`.
- **YAML-конфиг** — только если правило хорошо ложится на поле конфига;
  расширьте `internal/config` и `Config.ToRules()`. Если в YAML не ложится —
  оставляйте code-level rule.

### 3. Договорённости

- Доменная логика в `internal/...`, публичные пакеты остаются тонкими.
- По возможности точные сообщения и позиции в `Violation`.
- Тесты рядом с правилом (`internal/policy/*_test.go` или `internal/lint`).

---

## Загрузка конфига (CLI)

```text
FindAndLoad(startDir)  ->  поиск вверх .commitkit.yml
Load(path)             ->  явный файл
Default()              ->  если файла нет (Found=false)
Config.ToRules()       ->  lint.Rules с defaults CC для пропущенных полей
```

У CLI нет публичного пакета `config` — конфиг это application-деталь в
`internal/config`.

---

## Поверхность CLI

| Команда                           | Роль                                         |
| --------------------------------- | -------------------------------------------- |
| `check`                           | Parse + policy rules                         |
| `parse`                           | Печать структуры сообщения (debug / tooling) |
| `install-hook` / `uninstall-hook` | Управление `commit-msg` hook                 |
| `version` / `help`                | Мета                                         |

Форматирование вывода — в `cli/`; `internal/*` возвращает данные и ошибки.

---

## Лицензия

Контрибуции принимаются на условиях лицензии MIT.
