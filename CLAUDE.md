# Pokedex CLI

A CLI tool to explore the Pokemon world using the PokeAPI.

## Build & Run

```sh
make build        # compile binary
make run          # build and run
make test         # run all tests with -race -count=1
make lint         # run golangci-lint
make fmt          # format with gofmt + goimports
make check        # fmt + vet + lint + test (full pipeline)
```

## Project Structure

```
.
├── main.go                  # entry point — wiring and REPL loop only
├── repl.go                  # types (config, cliCommand), command registry, cleanInput
├── command_*.go             # one file per CLI command
├── internal/
│   ├── pokecache/
│   │   └── cache.go         # thread-safe in-memory cache with TTL
│   └── pokeapi/
│       ├── client.go        # HTTP client, API methods, response caching
│       └── types.go         # response structs
├── .claude/
│   ├── settings.json        # hooks (auto-format on edit)
│   └── skills/              # slash commands (/check, /test, /lint)
├── Makefile
├── .golangci.yml
└── CLAUDE.md
```

## Conventions

- **Go version**: 1.26+
- **Module**: `github.com/ahmed-abdelhamid/pokedexcli`
- **Binary**: `pokedexcli`
- **Linting**: golangci-lint with `.golangci.yml` — all code must pass `make lint`
- **Testing**: table-driven tests, use `github.com/google/go-cmp/cmp` for diffs
- **Formatting**: `gofmt` + `goimports` — tabs, no trailing whitespace
- **Error handling**: return errors up the call stack, handle at the boundary
- **Naming**: idiomatic Go — short, clear names; no stuttering (e.g. `config.Config` is bad)
- **Package comments**: every package needs a comment on the `package` declaration
- **Unexported by default**: only export what other packages need
- **One command per file**: each CLI command lives in its own file (e.g. `command_map.go`)
- **main.go stays thin**: only wiring and the REPL loop
- **No global state**: pass dependencies via structs/function parameters

## Workflow

Every change follows this sequence — no exceptions:

1. **Read first** — read CLAUDE.md, memory, and all files being modified before writing code
2. **Design first** — for non-trivial features, state the approach (file layout, struct design) before coding
3. **Small steps** — one logical change at a time; never mix refactors with features
4. **Verify after every change** — run `make check` after each meaningful edit; fix issues before moving on
5. **Consider tests** — if a function has logic worth testing, add a table-driven test
6. **Commit at logical boundaries** — commit after each complete, passing change (new feature, refactor, bugfix); atomic commits, clear messages, never skip hooks
7. **Update docs and memory** — after completing work, check if CLAUDE.md or memory need updating to reflect new structure, conventions, or learnings
8. **Evaluate tooling** — after completing work, consider if a new skill, hook, or command would prevent repetitive work in the future

## Slash Commands

- `/check` — run full pipeline (fmt + vet + lint + test), fix issues
- `/test` — run tests with race detection, investigate failures
- `/lint` — run golangci-lint, fix issues

## Hooks

- **PostToolUse (Write|Edit)** — auto-runs `goimports` on `.go` files after every edit
