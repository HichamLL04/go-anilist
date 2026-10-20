# GitHub Copilot Instructions — go-anilist

Please refer to [AGENTS.md](../AGENTS.md) for all architectural rules, coding guidelines, Go conventions, and repository invariants.

All code generation must follow:
- `CGO_ENABLED=0` compatibility.
- Zero external dependencies beyond the standard library.
- `context.Context` as the first parameter in any I/O operation.
- Explicit error wrapping with `%w`.
- Table-driven unit tests utilizing `net/http/httptest`.