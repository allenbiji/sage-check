# Architecture

This document describes Preboot's internal design — how packages are organized, how data flows through the system, and where the key extension points are.

---

## Package map

```
github.com/allenbiji/preboot/
├── cmd/preboot/          Entry point — wires up CLI and exits
├── internal/
│   ├── cli/           Cobra command definitions (root, check, init, validate)
│   ├── config/        YAML loading, merging, and validation
│   ├── detect/        Auto-detection of project frameworks (init command)
│   ├── engine/        Execution loop and result rendering
│   ├── checks/        Seven concrete check implementations
│   ├── model/         Shared data types (structs, enums)
│   └── registry/      Factory registry that maps type strings → constructors
```

All packages are `internal/` — the tool is a CLI binary, not a library.

---

## Data flow

### `preboot check` path

```
cmd/preboot/main.go
  └─ cli.Execute()
       └─ cli.NewCheckCmd()
            │
            ├─ config.Load()  ─────────────────────────────────────────┐
            │    ├─ reads preboot-auto.yml (if exists)                    │
            │    ├─ reads preboot.yml (if exists)                         │
            │    ├─ config.merge()  ← user overrides auto by name      │
            │    └─ config.MergeDefaults()  ← injects strict/timeout   │
            │                                                           │
            └─ engine.Run(cfg, quickMode) ◄─────────────────────────────┘
                 │
                 ├─ for each check in cfg.Checks:
                 │    ├─ (skip if quick && network check)
                 │    ├─ inject timeout_ms from defaults if missing
                 │    ├─ registry.Build(checkCfg)
                 │    │    └─ looks up Factory by checkCfg.Type
                 │    │         └─ factory(checkCfg) → Check
                 │    └─ check.Execute() → error | nil
                 │         └─ render result (icon + name + message + fix)
                 │
                 └─ print summary → return ErrCheckFailed | nil
```

### `preboot init` path

```
cmd/preboot/main.go
  └─ cli.Execute()
       └─ cli.NewInitCmd()
            └─ detect.ScanRepo()
                 ├─ detect.go()        ← go.mod → command_exists check
                 ├─ detect.makefile()  ← Makefile → command_exists check
                 ├─ detect.docker()    ← docker-compose.yml → port_free + command_exists
                 └─ detect.env()       ← .env.example → file_exists + env_exists checks
                      └─ detect.ExtractEnvKeys(path) → map[key]""
                           └─ write preboot-auto.yml
```

---

## Key interfaces and types

### `model.CheckConfig` — the unit of work

```go
// internal/model/config.go
type CheckConfig struct {
    Name     string
    Type     CheckType         // string alias: "command_exists", "port_free", ...
    Severity Severity          // string alias: "info", "warning", "blocker"
    Options  map[string]string // type-specific key/value pairs
    Message  string
    Fix      string
}
```

Every check definition in YAML becomes one `CheckConfig` value.

### `registry.Check` — the check interface

```go
// internal/registry/registry.go
type Check interface {
    Execute() error
}

type Factory func(cfg model.CheckConfig) (Check, error)
```

A `Factory` validates and constructs a `Check`; the `Check` runs the diagnostic.

### `registry` — the global factory map

```go
var factories = map[model.CheckType]Factory{}

func Register(t model.CheckType, f Factory) { factories[t] = f }
func Build(cfg model.CheckConfig) (Check, error) { ... }
func IsKnownType(t model.CheckType) bool { ... }
```

All seven check packages register themselves in their `init()` functions:

```go
// internal/checks/command.go
func init() {
    registry.Register(model.CommandExists, newCommandCheck)
}
```

The `cmd/preboot/main.go` entry point imports the `checks` package with a blank import to trigger these `init()` calls:

```go
import _ "github.com/allenbiji/preboot/internal/checks"
```

---

## Package responsibilities

### `cmd/preboot`

- Single file: `main.go`
- Calls `cli.Execute()`, which sets exit code and calls `os.Exit`
- The only code that calls `os.Exit`

### `internal/cli`

| File | Responsibility |
|---|---|
| `root.go` | Creates the root `preboot` Cobra command; attaches subcommands |
| `banner.go` | Prints the gradient ASCII art banner to stdout (terminal only) |
| `check.go` | `preboot check` — parses flags, calls `config.Load/LoadFrom`, calls `engine.Run` |
| `init.go` | `preboot init` — checks for existing file, calls `detect.ScanRepo`, writes YAML |
| `validate.go` | `preboot validate` — calls `config.Load/LoadFrom`, calls `config.ValidateConfig` |

### `internal/config`

| File | Responsibility |
|---|---|
| `load.go` | `Load()` and `LoadFrom()` — find, read, and unmarshal YAML; coordinate merging |
| `merge.go` | `mergeConfigs()` — combine auto + user configs by name; `MergeDefaults()` — inject global defaults into each check |
| `validate.go` | `ValidateConfig()` — enforce schema rules (version, name, severity, type) |

### `internal/detect`

| File | Responsibility |
|---|---|
| `repo.go` | `ScanRepo()` — orchestrates all detectors; returns `*model.PrebootConfig` |
| `go.go` | Detects `go.mod`; emits `go-installed` check |
| `docker.go` | Detects `docker-compose.yml`/`compose.yaml`; parses port mappings; emits checks |
| `env.go` | Detects `.env.example`/`.env.template`; emits checks for each key |
| `files.go` | `ExtractEnvKeys(path)` — generic `.env`-style file parser |

### `internal/engine`

| File | Responsibility |
|---|---|
| `run.go` | `Run(cfg, quick)` — the main execution loop; renders results; returns `ErrCheckFailed` |
| `color.go` | ANSI color helpers; `colorEnabled` bool (controls output in tests) |

### `internal/checks`

Seven files, one per check type. Each:
1. Defines a struct implementing `Check`
2. Implements `Execute() error`
3. Defines a factory function that validates options and constructs the struct
4. Registers itself in `init()`

| File | Type registered |
|---|---|
| `command.go` | `command_exists` |
| `file.go` | `file_exists` |
| `directory.go` | `directory_exists` |
| `env.go` | `env_exists` |
| `http.go` | `http_reachable` |
| `tcp.go` | `tcp_reachable` |
| `port.go` | `port_free` |

### `internal/model`

Pure data — no logic. Defines `PrebootConfig`, `CheckConfig`, `Severity`, and `CheckType` string aliases.

### `internal/registry`

Pure plumbing — no check logic. Owns the global factory map, `Register`, `Build`, `IsKnownType`.

---

## Color and output

The engine uses ANSI escape codes for colour. Color is enabled when all three conditions hold:

1. `stdout` is a terminal (`isatty.IsTerminal`)
2. `NO_COLOR` environment variable is not set
3. `TERM` environment variable is not `"dumb"`

This logic lives in `internal/engine/color.go`. Tests override the `colorEnabled` package-level variable to disable color in test output.

The banner in `internal/cli/banner.go` additionally gates on `isatty` and uses `charmbracelet/lipgloss` + `go-colorful` to render a green-to-teal gradient.

---

## Error handling conventions

| Error type | Represented as |
|---|---|
| Check failed (expected) | `Execute()` returns `error` with descriptive message |
| Unknown check type | `registry.Build()` returns `fmt.Errorf` — engine renders as internal error |
| Config parse/validation | `config.Load()` returns `error` — CLI prints to stderr, exits 2 |
| Internal engine error | `engine.Run()` returns `error` — CLI exits 2 |
| Blocker check failed | `engine.Run()` returns `engine.ErrCheckFailed` — CLI exits 1 |

---

## Dependency graph

```
cmd/preboot
  └── internal/cli
        ├── internal/config
        │     ├── internal/model
        │     └── internal/registry
        ├── internal/detect
        │     └── internal/model
        └── internal/engine
              ├── internal/model
              ├── internal/registry
              └── internal/checks  (via blank import in cmd/preboot)
                    ├── internal/registry
                    └── internal/model
```

`internal/model` and `internal/registry` are the two packages everyone imports. They have no imports of their own within the project — no circular dependency risk.

---

## Adding a new check type

See [Contributing — Adding a new check type](contributing.md#adding-a-new-check-type) for a step-by-step walkthrough.
