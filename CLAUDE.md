# CLAUDE.md — go-anilist

Project rules for Claude Code assistants:

Refer to [AGENTS.md](AGENTS.md) as the SINGLE SOURCE OF TRUTH for all repository
rules, Go style conventions, architectural invariants, resource limits, and prohibitions.

## Specific notes for Claude Code

- Never run tests that require real internet connectivity; all tests must use local mocks via `httptest.Server`.
- All modifications must follow Conventional Commits 1.0.0 (see [docs/ai/commits.md](docs/ai/commits.md)).
- Before proposing any external dependency, review [docs/design/graphql.md](docs/design/graphql.md).