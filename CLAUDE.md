# CLAUDE.md — go-anilist

Reglas de proyecto para asistentes Claude Code:

Consulta [AGENTS.md](AGENTS.md) como la FUENTE ÚNICA de verdad sobre convenciones,
estilo Go, diseño arquitectónico, límites de recursos y prohibiciones de este repositorio.

## Notas específicas para Claude Code

- No ejecutes tests que dependan de conexión a internet; todo test debe utilizar mocks con `httptest.Server`.
- Toda modificación debe documentarse y commitearse siguiendo Conventional Commits (ver [docs/ai/commits.md](docs/ai/commits.md)).
- Antes de proponer una nueva dependencia, revisa [docs/design/graphql.md](docs/design/graphql.md).
