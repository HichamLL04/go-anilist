# AI Agent Workflow — go-anilist

Methodology and rigorous workflow that all AI agents must follow when modifying this repository.

---

## 1. Core Principle: "Read Before Writing"

1. **Consult the Sources of Truth**:
   - `AGENTS.md` for project rules and invariants.
   - `docs/reference/anilist-api.md` to validate AniList official schema types and fields.
   - `docs/design/api-publica.md` to verify public API contracts.
2. **Do Not Invent APIs or Fields**:
   - Every field in a struct must match official AniList GraphQL v2 specifications.
   - Public signatures must follow established design decisions.

---

## 2. Managing Dependencies and Decisions

- **No Arbitrary Dependencies**: Any proposed addition to `go.mod` must be evaluated
  and formally justified in `docs/design/`.
- **Recording Doubts**: When facing ambiguities, document the open question in
  `docs/preguntas-abiertas.md` before making assumptions.

---

## 3. Feature Implementation Cycle

1. **Define the GraphQL Query / Mutation**: Document exact query strings and variables.
2. **Define Go Types**: Create compact, strongly typed Go structs.
3. **Write Tests with Fixtures**:
   - Add mock JSON responses in `testdata/`.
   - Write unit tests using `httptest.Server` verifying sent variables, headers, and output.
4. **Implement Service Methods**: Connect transport layer to the public service.
5. **Verify Quality**:
   - Run `go test -v -race ./...`.
   - Run `go vet ./...`.
6. **Atomic Commit**: Commit following Conventional Commits 1.0.0.