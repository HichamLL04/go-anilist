# Roadmap — `go-anilist`

Este documento detalla las fases de desarrollo planificadas para la librería cliente
`go-anilist`.

---

## Fase 0: Diseño Arquitectónico y Especificación
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
- [x] Inicialización del módulo Go (`go.mod`) con Go 1.27 y `CGO_ENABLED=0`.
- [x] Implementación de `Client`, constructor `NewClient` y opciones funcionales de configuración (`WithBaseURL`, `WithToken`, `WithUserAgent`, etc.).
- [x] Implementación del motor interno de transporte GraphQL (`doGraphQL`):
  - [x] Serialización de queries y variables JSON.
  - [x] Envío mediante `http.Client`.
  - [x] Deserialización en streaming con `json.NewDecoder`.
  - [x] Captura y extracción de cabeceras `X-RateLimit-*`.
  - [x] Normalización de errores estructurados en `ErrorResponse` y `GraphQLError`.
- [x] Suite de pruebas base con `httptest.Server` para transporte, cancelaciones de contexto y respuestas 429.

---

## Fase 2: Servicios de Usuario y Medios (`User` & `Media`)
- [x] Implementación de `UserService`:
  - [x] Consulta `Viewer` para obtener ID de usuario, nombre y formato de puntuación.
  - [x] Tests unitarios con fixtures JSON en `testdata/viewer_success.json`.
- [x] Implementación de `MediaService`:
  - [x] Búsqueda de anime por título (`Search`) con opciones de formato y año.
  - [x] Consulta directa por ID de AniList (`GetByID`) y por ID de MyAnimeList (`GetByIDMal`).
  - [x] Tests unitarios con fixtures JSON en `testdata/media_search_success.json`.

---

## Fase 3: Gestión de Listas y Progreso (`MediaList`)
- [x] Implementación de `MediaListService`:
  - [x] Consulta de entrada individual (`GetEntry`) por `userId` y `mediaId`.
  - [x] Consulta de colección paginada (`GetCollection`) con filtro por estado.
  - [x] Mutación `SaveEntry` para actualizar episodio visto (`progress`), estado (`status`) y nota (`score`).
  - [x] Mutación `DeleteEntry` para eliminar entradas.
- [x] Suite exhaustiva de pruebas unitarias para todas las variantes de mutación de listas.

---

## Fase 4: Autenticación OAuth2
- [x] Implementación de `AuthService`:
  - [x] Generador de URL de autorización `AuthorizeURL`.
  - [x] Intercambio de código `ExchangeCode` contra el endpoint `/oauth/token`.
  - [x] Renovación automática mediante `RefreshToken`.
- [x] Tests unitarios simulando endpoints OAuth2 locales.

---

## Fase 5: Estabilización y Tagging Inicial
- [x] Validación completa de tests con `go test -v ./...`.
- [x] Verificación de código estático con `go vet ./...`.
- [ ] Etiquetado del release inicial `v0.1.0`.