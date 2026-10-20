# GraphQL Integration Strategy — `go-anilist`

This document compares three approaches for interacting with the AniList GraphQL
API in Go, evaluating dependencies, binary footprint, maintenance, and allocations.

---

## 1. Evaluated Technical Approaches

### Option A: Handcrafted Queries with Typed Structs (`net/http` + `encoding/json`)
- **Description**: Queries and mutations are declared as Go string constants. Responses are streamed directly into typed Go structs via `json.NewDecoder`.
- **Dependencies**: **0**. Standard library only (`net/http`, `encoding/json`, `context`, `fmt`).

### Option B: Code Generation (e.g. `genqlient`)
- **Description**: Requires downloading the entire AniList SDL schema (`schema.graphql`), maintaining `.graphql` files, and generating Go structs via code generation.
- **Dependencies**: Build tooling dependency (`genqlient`), `genqlient.yaml`, and extra `go generate` build steps.

### Option C: Generic Dynamic GraphQL Clients (e.g. `machinebox/graphql`)
- **Description**: Third-party packages encapsulating GraphQL queries via reflection at runtime.
- **Dependencies**: 1+ external dependencies in `go.mod`.

---

## 2. Comparative Matrix

| Criterion | Option A: Handcrafted Queries | Option B: Generator (`genqlient`) | Option C: Generic Clients |
|-----------|:-----------------------------:|:----------------------------------:|:-------------------------:|
| **`go.mod` Dependencies** | **0 (Stdlib only)** **[DOC-GO]** | 1 tooling dependency (`genqlient`) | 1+ external libraries |
| **Binary Footprint [RE]** | **Minimal (~30-50 KiB compiled)** | Moderate (~250-400 KiB generated) | Medium (~300-600 KiB reflection) |
| **Build Step** | Instant (`go build ./...`) | Requires `go generate` + SDL sync | Instant (`go build ./...`) |
| **Schema Maintenance** | Only fields actually consumed are defined | Requires syncing 2-3 MB SDL schema | Dynamic runtime maintenance |
| **Memory Allocations** | Optimal (streaming with `json.NewDecoder`) | Very good | Higher due to reflection |
| **Network Control** | Direct byte-by-byte header & payload control | Indirect | Rigid HTTP abstractions |

---

## 3. Decision and Recommendation

**Option A (Handcrafted queries with typed structs) is adopted.**

### Justification:
1. **Zero External Dependencies**: Guarantees maximum long-term stability and security.
2. **Maximum Efficiency**: `json.NewDecoder` decodes directly from the response stream, avoiding buffer duplicates in memory.
3. **Simplicity**: The scope consists of exactly 4 queries and 2 mutations; the overhead of code generation is unjustified.
4. **Rate Limit Inspection**: Allows natural extraction of `X-RateLimit-*` headers on `*http.Response`.