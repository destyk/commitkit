# commitkit

**A small Go-native tool that checks whether your Git commit messages follow the rules you set.**

commitkit turns a raw commit message into a structured model (type, scope,
breaking marker, description, body, Git trailers) and validates it with a
small rule engine. **Conventional Commits is implemented as a policy**, not
hard-wired into the parser — so the same parser stays useful when a project
uses a different convention.

[Русская версия](README.ru.md) · [Contributing](CONTRIBUTING.md)

---

## Problems it solves

| Problem                                 | How commitkit helps                                                       |
| --------------------------------------- | ------------------------------------------------------------------------- |
| Commit messages drift across the team   | One config (`.commitkit.yml`) or a shared policy in code                  |
| Tools hard-code Conventional Commits    | Parser is convention-agnostic; CC is just the default policy              |
| Diagnostics are vague                   | Source positions on every part of the message                             |
| Git trailers are messy                  | Proper handling of continuations, multiline values, `BREAKING CHANGE`     |
| Hooks and CI need different tools       | Same CLI for local `commit-msg` hooks and pipelines                       |
| Library users pull in a heavy framework | Thin public adapters (`commit`, `lint`, `policy`); core under `internal/` |

---

## Quick start

### Install

```bash
# Install
# Make sure `$(go env GOPATH)/bin` is on your `PATH`.
go install github.com/destyk/commitkit@latest

# or from a local clone
git clone https://github.com/destyk/commitkit.git
cd commitkit
make build
```

### Prebuilt binaries

Download a release binary for your OS/arch from:

**https://github.com/destyk/commitkit/releases**

Unpack it, place `commitkit` on your `PATH`, and verify:

```bash
commitkit version
```

### Usage

```bash
# Optional project config (walk-up from cwd). No file -> Conventional Commits defaults.
cat > .commitkit.yml <<'YAML'
types: [feat, fix, docs, chore]
description:
  min: 1
  max: 72
  lowercase: true
scope:
  required: false
YAML

# Check the latest commit
git log -1 --format=%B | commitkit check

# Or a file (e.g. during a hook)
commitkit check --file .git/COMMIT_EDITMSG

# Install a git commit-msg hook
commitkit install-hook
```

---

## Main commands

```bash
commitkit check [--file FILE] [--config FILE]
commitkit parse [--file FILE]
commitkit install-hook [--force] [--path DIR]
commitkit uninstall-hook [--force] [--path DIR]
commitkit version
```

| Exit | Meaning                     |
| ---- | --------------------------- |
| 0    | Success                     |
| 1    | Parse / policy / hook error |
| 2    | Usage / I/O / config error  |

---

## Configuration

`.commitkit.yml` is searched upward from the current directory. Use `--config`
to pass an explicit path.

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

Omitted fields fall back to Conventional Commits defaults.

---

## Why this design

- **Policy, not hard-coding** — swap or compose rules without rewriting the parser.
- **Precise diagnostics** — positions on type, scope, description, trailers.
- **CLI vs library** — console output stays in the CLI; domain packages return structures.
- **Small surface** — import `commit`, `lint`, `policy` only; never `internal/...`.

---

## Contributing / internals

Architecture, public adapters, how to add a rule, and library examples are in
**[CONTRIBUTING.md](CONTRIBUTING.md)**.

Start there if you want to extend commitkit or embed it in your own tools.

---

## License

MIT
