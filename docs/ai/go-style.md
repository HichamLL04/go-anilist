# Reglas de Estilo Go — go-anilist

Guía de estilo y directrices técnicas obligatorias para cualquier código Go
en la librería `go-anilist`.

---

## 1. Manejo de errores

- **Envolver siempre con `%w`**: Todos los errores propagados deben envolverse usando
  `fmt.Errorf("anilist: <op>: %w", err)`. Nunca descartar errores con `_` en rutas
  de ejecución crítica.
- **Errores tipados y centinela**:
  - Definir errores centinela para condiciones previsibles (ej. `ErrNotFound`, `ErrRateLimited`, `ErrUnauthorized`).
  - Proporcionar un tipo estructurado `GraphQLError` y `ErrorResponse` que exponga los códigos de estado HTTP, mensajes devueltos por AniList y ubicaciones en la consulta GraphQL.
  - Implementar métodos `errors.Is` / `errors.As` para permitir que el consumidor inspeccione la causa subyacente.

---

## 2. Contexto y cancelación

- **`context.Context` como primer parámetro**: En cualquier método o función que realice peticiones HTTP, I/O o cálculos asíncronos:
  ```go
  func (s *MediaService) Get(ctx context.Context, id int) (*Media, *Response, error)
  ```
- **Respetar cancelación**: Verificar `ctx.Err()` en operaciones en bucle y propagar el contexto a `http.Request.WithContext(ctx)`.
- Si `ctx == nil`, rechazar la petición con un error explícito antes de emitir tráfico de red.

---

## 3. Estado, variables globales y funciones `init()`

- **Prohibidas las variables globales mutables**: Cero estado mutable compartido. Cada instancia de `Client` gestiona su propia configuración, transporte y caché opcional.
- **Sin `init()`**: No realizar registros mágicos en `init()`. La configuración debe ser explícita mediante opciones en `NewClient`.
- Las constantes de endpoints y URLs deben ser inmutables.

---

## 4. Diseño de interfaces y API pública

- **Interfaces pequeñas y orientadas al consumidor**:
  - `go-anilist` expone structs concretos para sus servicios (`MediaService`, `MediaListService`, etc.) y tipos de datos devueltos.
  - No forzar interfaces gigantes en la librería; las interfaces se definen en el paquete consumidor (ej. en `aniSync` mediante la interfaz `Tracker`).
- **Opciones funcionales tipadas**: Para parámetros opcionales y paginación, usar el patrón de opciones funcionales tipadas (ej. `Page(1)`, `PerPage(50)`, `Sort(MediaSortTitle)`).
- **Inyección de dependencias**:
  - El cliente debe aceptar un `*http.Client` preconfigurado:
    ```go
    func NewClient(httpClient *http.Client, opts ...ClientOption) *Client
    ```
  - Si `httpClient` es `nil`, inicializar uno por defecto con timeouts razonables.

---

## 5. Estrategia de testing

- **Tests table-driven**: Uso de tablas de tests con subtests nombrados mediante `t.Run(tc.name, func(t *testing.T) { ... })`.
- **Cero llamadas de red reales**: Los tests unitarios jamás deben depender de los servidores de AniList. Se debe emplear exclusivamente `net/http/httptest.Server` con fixtures JSON precargadas.
- **Race detector**: Los tests deben pasar limpios con `go test -race ./...`.

---

## 6. Build y compilación

- **`CGO_ENABLED=0`**: Todo el código debe compilar sin soporte CGO para garantizar binarios estáticos portables.
- **Dependencias externas**:
  - Minimizar dependencias. Preferir siempre la biblioteca estándar (`net/http`, `encoding/json`, `context`, `time`).
  - No añadir clientes GraphQL pesados o librerías de terceros sin justificación documentada en `docs/design/graphql.md`.
