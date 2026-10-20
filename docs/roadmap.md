# Roadmap — `go-anilist`

Este documento detalla las fases de desarrollo planificadas para la librería cliente
`go-anilist`.

---

## Fase 0: Diseño Arquitectónico y Especificación (Actual)
- [x] Reglas estrictas de desarrollo y configuración para agentes de IA (`AGENTS.md`, `docs/ai/*`).
- [x] Análisis arquitectónico comparativo de la librería de referencia `go-myanimelist` (`docs/reference/go-myanimelist.md`).
- [x] Documentación y especificación detallada de la API GraphQL de AniList (`docs/reference/anilist-api.md`).
- [x] Delimitación estricta del alcance funcional necesario para aniSync (`docs/design/alcance.md`).
- [x] Comparativa de tecnologías GraphQL y decisión por cliente ligero sin dependencias (`docs/design/graphql.md`).
- [x] Especificación formal de tipos, interfaces y métodos de la API pública (`docs/design/api-publica.md`).
- [x] Definición de la estrategia de pruebas unitarias deterministas con `httptest` (`docs/design/testing.md`).
- [x] Establecimiento de la política SemVer y relación con Conventional Commits (`docs/design/versionado.md`).

---

## Fase 1: Esqueleto del Paquete y Capa de Transporte
- [ ] Inicialización del módulo Go (`go.mod`) con Go 1.25.x / 1.24.x y `CGO_ENABLED=0`.
- [ ] Implementación de `Client`, constructor `NewClient` y opciones funcionales de configuración.
- [ ] Implementación del motor interno de transporte GraphQL (`doGraphQL`):
  - Serialización de queries y variables JSON.
  - Envío mediante `http.Client`.
  - Deserialización en streaming con `json.NewDecoder`.
  - Captura y extracción de cabeceras `X-RateLimit-*`.
  - Normalización de errores estructurados en `ErrorResponse` y `GraphQLError`.
- [ ] Suite de pruebas base con `httptest.Server` para transporte, cancelaciones de contexto y respuestas 429.

---

## Fase 2: Servicios de Usuario y Medios (`User` & `Media`)
- [ ] Implementación de `UserService`:
  - Consulta `Viewer` para obtener ID de usuario, nombre y formato de puntuación.
  - Tests unitarios con fixtures JSON en `testdata/viewer_success.json`.
- [ ] Implementación de `MediaService`:
  - Búsqueda de anime por título (`Search`) con opciones de formato y año.
  - Consulta directa por ID de AniList (`GetByID`) y por ID de MyAnimeList (`GetByIDMal`).
  - Tests unitarios con fixtures JSON en `testdata/media_search_success.json`.

---

## Fase 3: Gestión de Listas y Progreso (`MediaList`)
- [ ] Implementación de `MediaListService`:
  - Consulta de entrada individual (`GetEntry`) por `userId` y `mediaId`.
  - Consulta de colección paginada (`GetCollection`) con filtro por estado.
  - Mutación `SaveEntry` para actualizar episodio visto (`progress`), estado (`status`) y nota (`score`).
  - Mutación `DeleteEntry` para eliminar entradas.
- [ ] Suite exhaustiva de pruebas unitarias para todas las variantes de mutación de listas.

---

## Fase 4: Autenticación OAuth2
- [ ] Implementación de `AuthService`:
  - Generador de URL de autorización `AuthorizeURL`.
  - Intercambio de código `ExchangeCode` contra el endpoint `/oauth/token`.
  - Renovación automática mediante `RefreshToken`.
- [ ] Tests unitarios simulando endpoints OAuth2 locales.

---

## Fase 5: Optimización, Benchmarks y Release v0.1.0
- [ ] Benchmarks de allocations y consumo de memoria (`testing.B`).
- [ ] Verificación de linters (`golangci-lint run`).
- [ ] Integración inicial con `aniSync` para validar la interfaz `Tracker`.
- [ ] Publicación del tag `v0.1.0`.
