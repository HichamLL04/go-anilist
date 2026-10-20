# AGENTS.md — go-anilist

Fuente única de reglas para agentes de IA que trabajen en este repositorio.
Lee este fichero completo antes de tocar cualquier código.

---

## 1. Propósito del proyecto

**go-anilist** es una librería Go cliente para la API GraphQL (v2) de AniList.
Su objetivo es proporcionar un cliente tipado, idiomático, eficiente en memoria
y con cero dependencias pesadas, para ser consumido primariamente por **aniSync**
y por cualquier aplicación Go que requiera interactuar con AniList.

La librería no implementa lógica de dominio de aniSync ni maneja webhooks de
Plex; es una biblioteca cliente agnóstica e independiente.

---

## 2. Mapa de `docs/`

```
docs/
  ai/                        ← reglas para agentes IA (este árbol)
    commits.md               ← convención Conventional Commits + ejemplos
    go-style.md              ← reglas de estilo Go
    arquitectura-reglas.md   ← invariantes arquitectónicos de la librería
    presupuesto-recursos.md  ← reglas de eficiencia y memoria
    flujo-de-trabajo.md      ← flujo de trabajo paso a paso
  reference/
    go-myanimelist.md        ← análisis de arquitectura de go-myanimelist
    anilist-api.md           ← especificación y detalles de AniList GraphQL API
  design/
    alcance.md               ← alcance del subconjunto de AniList soportado
    graphql.md               ← comparativa de enfoques GraphQL y justificación
    api-publica.md           ← firmas y diseño de la API Go pública
    testing.md               ← estrategia de testing con httptest y fixtures
    versionado.md            ← política de SemVer y ciclo de releases
  roadmap.md                 ← fases de implementación
  preguntas-abiertas.md      ← dudas pendientes de resolución
```

---

## 3. Comandos de build / test / lint

> **TODO**: definir cuando exista código Go.

```bash
# Compilar (cuando exista go.mod)
# CGO_ENABLED=0 go build ./...

# Tests unitarios (con race detector)
# go test -race ./...

# Lint estricto
# golangci-lint run
```

---

## 4. Reglas de estilo Go

Ver detalles en `docs/ai/go-style.md`. Resumen:

- Errores siempre envueltos con `%w` (`fmt.Errorf("anilist: %s: %w", op, err)`).
- `context.Context` como primer parámetro en toda función que realice I/O o HTTP.
- Sin variables globales mutables ni `init()`.
- Cliente configurable mediante constructor `NewClient(httpClient *http.Client, opts ...Option)`.
- Opciones funcionales tipadas para filtros y mutaciones.
- Tests table-driven con `t.Run` y uso exclusivo de `httptest.Server` (sin red en tests).
- `CGO_ENABLED=0` estricto.
- Prohibido introducir dependencias externas pesadas o clientes GraphQL de terceros sin evaluación justificada.

---

## 5. Convención de commits

Ver detalles y ejemplos en `docs/ai/commits.md`.

Formato: `<type>(<scope>)<!>: <description>`

- Tipos: `feat fix docs style refactor perf test build ci chore revert`
- Scope obligatorio cuando aplique (ej. `docs(ref)`, `docs(design)`, `feat(client)`).
- Descripción en inglés, imperativo, minúscula, sin punto, ≤ 72 chars.
- Un commit por unidad lógica atómica. Nunca `--amend` ni `push --force`.
- Nunca commitear secretos, tokens ni credenciales personales.

---

## 6. Presupuesto de recursos y memoria

Ver `docs/ai/presupuesto-recursos.md`.

- Librería ultraligera: sin dependencias externas innecesarias (`stdlib` + opcional rate limiter si se justifica).
- Decodificación directa con `json.NewDecoder` hacia structs compactos.
- No duplicar buffers ni retener respuestas completas en memoria si no se requiere.
- Paginación estricta (máximo 50 elementos por página según API de AniList).

---

## 7. Etiquetas de procedencia

Toda afirmación técnica en la documentación debe indicar su origen:

| Etiqueta | Significado |
|----------|-------------|
| **[INFERIDO]** | Deducido del código fuente, tipos o contexto de uso |
| **[RE]** | Medido o analizado con herramientas de ingeniería inversa / profiling |
| **[DOC-GO]** | Extraído de la documentación oficial de Go o librerías estándar |
| **[DOC-ANILIST]** | Extraído de la documentación oficial de AniList GraphQL v2 |

---

## 8. Cosas PROHIBIDAS

- Acoplar `go-anilist` a conceptos de `aniSync` (Plex, webhooks, base de datos local).
- Importar paquetes o lógica de `aniSync` en `go-anilist`.
- Usar variables globales mutables o estado estático compartido.
- Realizar peticiones HTTP reales a AniList en suites de tests (usar `httptest.Server`).
- Hardcodear tokens OAuth2, client IDs, client secrets o credenciales de usuario.
- Añadir dependencias externas sin justificación formal en `docs/design/`.
- Ignorar o descartar la cancelación de `context.Context`.
- Reescribir la historia de Git (`git commit --amend`, `git push --force`).

---

## 9. Checklist "antes de dar una tarea por terminada"

- [ ] La documentación afectada está actualizada y sincronizada.
- [ ] El código compila sin errores cuando aplique (`go build ./...`).
- [ ] Los tests unitarios pasan sin red (`go test -race ./...`).
- [ ] Todos los errores devueltos envuelven la causa con `%w`.
- [ ] `context.Context` es el primer parámetro en llamadas con I/O.
- [ ] No hay variables globales mutables ni dependencias no autorizadas.
- [ ] No existen credenciales, tokens ni datos personales en el commit.
- [ ] El commit sigue rigurosamente Conventional Commits 1.0.0.
