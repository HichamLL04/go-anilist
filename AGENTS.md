# AGENTS.md — go-anilist

Single source of truth for AI agents working in this repository.
Read this entire file before modifying any code or documentation.

---

## 1. Project Purpose

**go-anilist** is an idiomatic Go client library for the AniList GraphQL API (v2).
Its goal is to provide a strongly typed, memory-efficient client with zero external
dependencies, primarily consumed by **aniSync** and any Go application that requires
interaction with AniList.

This library does not implement aniSync domain logic, Plex webhook handling, or local
database operations; it is a standalone, agnostic client library.

---

## 2. Map of `docs/`

```
docs/
  ai/                        ← Rules for AI agents
    commits.md               ← Conventional Commits 1.0.0 rules and examples
    go-style.md              ← Go style guide and error wrapping rules
    arquitectura-reglas.md   ← Library architectural invariants
    presupuesto-recursos.md  ← Resource efficiency and memory budget rules
    flujo-de-trabajo.md      ← Step-by-step agent workflow
  reference/
    go-myanimelist.md        ← Architectural analysis of go-myanimelist
    anilist-api.md           ← AniList GraphQL API technical specifications
  design/
    alcance.md               ← Scope and bounded context of supported features
    graphql.md               ← Evaluation of GraphQL integration approaches
    api-publica.md           ← Exact Go public API contracts and signatures
    testing.md               ← Deterministic testing strategy with httptest
    versionado.md            ← SemVer policy and release lifecycle
  roadmap.md                 ← Implementation roadmap and phases
  preguntas-abiertas.md      ← Open architectural questions and pending decisions
```

---

## 3. Build / Test / Lint Commands

```bash
# Compile library
CGO_ENABLED=0 go build ./...

# Run unit tests with race detection
go test -v -race ./...

# Static analysis
go vet ./...
```

---

## 4. Go Style Rules

See full details in `docs/ai/go-style.md`. Summary:

- Errors must always be wrapped using `%w` (`fmt.Errorf("anilist: %s: %w", op, err)`).
- `context.Context` must be the first parameter in all functions performing I/O or network calls.
- No mutable global variables and no `init()` functions.
- Client configured via functional options in `NewClient(httpClient *http.Client, opts ...ClientOption)`.
- Table-driven unit tests using `t.Run` and exclusive use of `httptest.Server` (zero real-network calls in tests).
- Strict `CGO_ENABLED=0` compatibility.
- Zero external dependencies beyond the Go standard library.

---

## 5. Commit Convention

See full details and examples in `docs/ai/commits.md`.

Format: `<type>(<scope>)<!>: <description>`

- Allowed types: `feat fix docs style refactor perf test build ci chore revert`
- Mandatory scope when applicable (e.g. `feat(client)`, `docs(design)`).
- Description in English, imperative mode, lowercase, no trailing period, ≤ 72 characters.
- Atomic commits per logical change. Never `--amend` or `push --force`.
- Never commit secrets, tokens, or personal credentials.

---

## 6. Resource and Memory Budget

See `docs/ai/presupuesto-recursos.md`.

- Ultra-lightweight: zero external dependencies (`net/http`, `encoding/json`, `context`).
- Streaming deserialization with `json.NewDecoder` directly on response bodies.
- Do not duplicate buffers or retain raw HTTP responses in memory unnecessarily.
- Strict pagination awareness (maximum 50 items per page enforced by AniList).

---

## 7. Provenance Tags

Technical assertions in documentation must specify their origin:

| Tag | Meaning |
|-----|---------|
| **[INFERRED]** | Deduced from source code, types, or usage context |
| **[RE]** | Measured or reverse-engineered via binary profiling |
| **[DOC-GO]** | Extracted from official Go documentation or standard library docs |
| **[DOC-ANILIST]** | Extracted from official AniList GraphQL v2 documentation |

---

## 8. Prohibited Actions

- Coupling `go-anilist` to `aniSync` concepts (Plex, webhooks, SQLite, sync loops).
- Importing `aniSync` packages or logic into `go-anilist`.
- Using mutable global variables or shared static state.
- Making real HTTP requests to AniList during unit tests (use `httptest.Server`).
- Hardcoding OAuth2 tokens, client IDs, secrets, or personal credentials.
- Adding third-party external dependencies without formal justification in `docs/design/`.
- Ignoring or discarding `context.Context` cancellation.
- Rewriting Git history (`git commit --amend`, `git push --force`) without explicit request.

---

## 9. Definition of Done Checklist

- [ ] All affected documentation is up to date.
- [ ] Code compiles without errors (`go build ./...`).
- [ ] Unit tests pass cleanly with race detector (`go test -v -race ./...`).
- [ ] `go vet ./...` reports 0 issues.
- [ ] All returned errors wrap the underlying cause with `%w`.
- [ ] `context.Context` is the first parameter in I/O operations.
- [ ] No mutable global variables or unauthorized dependencies exist.
- [ ] No secrets, tokens, or personal credentials are included.
- [ ] Commits strictly follow Conventional Commits 1.0.0.