# Commit Convention — go-anilist

All commits in this repository must strictly follow the **Conventional Commits 1.0.0** specification.

---

## 1. Commit Structure

```
<type>(<scope>)<!>: <description>

[optional body: lines ≤ 72 characters]

[optional footer(s): BREAKING CHANGE:, Refs:]
```

### Mandatory Rules
- **type**: One of the allowed types below (in lowercase).
- **scope**: Mandatory whenever the change targets a specific subsystem.
- **description**: In English, imperative mood ("add", "fix", "update", never "added" or "adds"), lowercase start, no trailing period, maximum 72 characters.
- **breaking change**: Append `!` immediately after type/scope (e.g. `feat(client)!: change auth signature`) and include `BREAKING CHANGE: <desc>` in the footer.
- **atomicity**: One commit per logical unit of work. Never bundle unrelated changes into a single monolithic commit.

---

## 2. Allowed Types

| Type | Purpose | Example |
|------|---------|---------|
| `feat` | New public feature or method | Implement MediaListService |
| `fix` | Bug fix or error handling correction | Fix GraphQL error array parsing |
| `docs` | Documentation-only changes | Add API specification |
| `style` | Formatting, whitespace (no code change) | Run gofmt on source files |
| `refactor` | Refactoring code without altering public API | Simplify internal JSON decoder |
| `perf` | Performance or memory allocation optimization | Reuse byte buffers in encoder |
| `test` | Adding or updating unit tests | Add httptest search fixtures |
| `build` | Build system or tooling dependencies | Adjust Go toolchain version |
| `ci` | Continuous integration workflows | Add GitHub Actions test matrix |
| `chore` | Maintenance tasks | Initialize .gitignore or licenses |
| `revert` | Reverting a previous commit | Revert previous commit |

---

## 3. Common Scopes in go-anilist

- `client`: Base HTTP / GraphQL client core.
- `media`: Anime/media lookup and search service.
- `medialist`: User tracking list query and mutation service.
- `user`: User profile and Viewer service.
- `auth`: OAuth2 token exchange and refresh service.
- `graphql`: Internal GraphQL transport, payload encoding, error handling.
- `design`: Architecture and design documentation.
- `reference`: API specifications and reference analysis.
- `ai`: Guidelines for AI agents.
- `repo`: Repository metadata and configuration.

---

## 4. Commit Examples

### Good Examples

```text
docs(reference): document go-myanimelist design analysis and patterns

Analyze package structure, service delegation, functional options,
and testing techniques for comparison with AniList GraphQL needs.
```

```text
feat(medialist): add UpdateMediaListEntry method and options

Implement SaveMediaListEntry mutation wrapper supporting score,
status, progress and repeat increments.
```

```text
fix(graphql): handle GraphQL error arrays with missing HTTP status

AniList returns HTTP 200 with top-level "errors" array when schema
validation or field resolution fails.
```

### Prohibited Examples (BAD)

- ❌ `added some files`: Not Conventional Commits; past tense.
- ❌ `WIP`: Incomplete, non-descriptive message.
- ❌ `docs: changes.`: Missing scope when applicable; trailing period.
- ❌ `feat(client): Add new method and also update readme and fix bug`: Not atomic; combines 3 separate concerns.
- ❌ `feat(client): implement huge client refactoring with all services`: Giant monolithic commit.