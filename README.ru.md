# commitkit

**Маленький инструмент на Go, который проверяет, что сообщения Git-коммитов
следуют правилам, которые задаёте вы.**

commitkit превращает сырое сообщение коммита в структурированную модель
(type, scope, breaking-маркер, description, body, Git trailers) и проверяет
её небольшим движком правил. **Conventional Commits — это политика**, а не
жёстко вшитая логика парсера: тот же парсер остаётся полезным при другой
конвенции.

[English version](README.md) · [Contributing](CONTRIBUTING.ru.md)

---

## Какие проблемы закрывает

| Проблема                                         | Как помогает commitkit                                                     |
| ------------------------------------------------ | -------------------------------------------------------------------------- |
| Сообщения коммитов разъезжаются по команде       | Один конфиг (`.commitkit.yml`) или общая policy в коде                     |
| Инструменты жёстко зашивают Conventional Commits | Парсер не привязан к конвенции; CC — политика по умолчанию                 |
| Диагностика размытая                             | Позиции в исходнике на каждую часть сообщения                              |
| Git trailers ведут себя непредсказуемо           | Continuations, multiline values, `BREAKING CHANGE`                         |
| Для хуков и CI нужны разные тулы                 | Один CLI для локального `commit-msg` и пайплайнов                          |
| Библиотека тащит тяжёлый фреймворк               | Тонкие публичные адаптеры (`commit`, `lint`, `policy`); ядро в `internal/` |

---

## Быстрый старт

### Установка

```bash
# Установка
# Убедитесь, что `$(go env GOPATH)/bin` есть в `PATH`.
go install github.com/destyk/commitkit@latest

# или клонируйте репозиторий
git clone https://github.com/destyk/commitkit.git
cd commitkit
make build
```

### Готовые бинарники

Скачайте бинарник под свою ОС/архитектуру со страницы релизов:

**https://github.com/destyk/commitkit/releases**

Распакуйте, положите `commitkit` в `PATH` и проверьте:

```bash
commitkit version
```

### Использование

```bash
# Опциональный конфиг проекта (поиск вверх от cwd). Нет файла -> defaults CC.
cat > .commitkit.yml <<'YAML'
types: [feat, fix, docs, chore]
description:
  min: 1
  max: 72
  lowercase: true
scope:
  required: false
YAML

# Проверить последний коммит
git log -1 --format=%B | commitkit check

# Или файл (например в хуке)
commitkit check --file .git/COMMIT_EDITMSG

# Установить git commit-msg hook
commitkit install-hook
```

---

## Основные команды

```bash
commitkit check [--file FILE] [--config FILE]
commitkit parse [--file FILE]
commitkit install-hook [--force] [--path DIR]
commitkit uninstall-hook [--force] [--path DIR]
commitkit version
```

| Код | Значение                     |
| --- | ---------------------------- |
| 0   | Успех                        |
| 1   | Ошибка parse / policy / hook |
| 2   | Usage / I/O / config         |

---

## Конфигурация

`.commitkit.yml` ищется вверх от текущей директории. Явный путь — через
`--config`.

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

Отсутствующие поля берутся из defaults Conventional Commits.

---

## Почему так устроено

- **Policy, а не hard-code** — правила можно менять и собирать, не переписывая парсер.
- **Точная диагностика** — позиции на type, scope, description, trailers.
- **CLI vs библиотека** — вывод в консоль в CLI, domain-пакеты возвращают структуры.
- **Маленькая поверхность** — импортируйте только `commit`, `lint`, `policy`.

---

## Участие в разработке / внутренности

Архитектура, публичные адаптеры, как добавить правило и примеры library API —
в **[CONTRIBUTING.ru.md](CONTRIBUTING.ru.md)**.

Начните оттуда, если хотите расширять commitkit или встраивать его в свои инструменты.

---

## Лицензия

MIT
