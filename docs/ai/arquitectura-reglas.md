# Invariantes Arquitectónicos — go-anilist

Reglas e invariantes estructurales que deben respetarse en todo momento
en el diseño e implementación de `go-anilist`.

---

## 1. Independencia y desacoplamiento

- **Totalmente agnóstica de aniSync**: `go-anilist` es una biblioteca autónoma.
  Bajo ningún concepto debe importar paquetes de `aniSync`, hacer referencia a
  Plex, SQLite, lógica de matching de biblioteca o webhooks.
- **Sin estado global**: Múltiples instancias de `Client` deben poder coexistir
  de forma concurrente y segura dentro del mismo proceso Go sin interferir
  entre sí.

---

## 2. Separación de capas: Transporte GraphQL vs API Pública Tipada

1. **Capa de Transporte Interno (`internal/transport` o método privado `doGraphQL`)**:
   - Encargada de serializar la petición GraphQL (`query`, `variables`).
   - Envía la solicitud HTTP POST hacia `https://graphql.anilist.co`.
   - Procesa los headers de rate limit (`X-RateLimit-*`).
   - Deserializa el sobre estándar GraphQL:
     ```json
     {
       "data": { ... },
       "errors": [ ... ]
     }
     ```
   - Normaliza errores HTTP y errores GraphQL en tipos Go enriquecidos.

2. **Capa de Servicios Públicos (`MediaService`, `MediaListService`, `UserService`, `AuthService`)**:
   - Exponen métodos de alto nivel con nombres de dominio idiomáticos en Go:
     - `media.GetByID(ctx, id)`
     - `media.Search(ctx, query, opts...)`
     - `mediaList.Get(ctx, userID, mediaID)`
     - `mediaList.Save(ctx, entry, opts...)`
   - Ocultan completamente los detalles sintácticos de las cadenas GraphQL y las mutaciones.
   - Retornan modelos Go fuertemente tipados (`*Media`, `*MediaListEntry`, `*User`).

---

## 3. Manejo de Rate Limiting y Cuotas

- **Lectura estricta de headers**: La capa de transporte debe registrar los headers
  `X-RateLimit-Limit`, `X-RateLimit-Remaining` y `X-RateLimit-Reset` en cada respuesta.
- **Soporte de `Response` metadata**: Todo método público retorna `(*T, *Response, error)`
  donde `Response` incluye los metadatos de rate limit para que el llamador (aniSync)
  pueda planificar backoff o throttling inteligente.
- **Respeto a HTTP 429**: Cuando la API devuelve 429 Too Many Requests, se debe extraer
  el header `Retry-After` y retornar un error estructurado `RateLimitError`.

---

## 4. Inyección y ciclo de vida de `http.Client`

- El cliente jamás creará instancias de transporte subyacentes que queden huérfanas.
- Se debe soportar la inyección de `http.Client` personalizados, permitiendo:
  - Autenticación OAuth2 delegada (ej. `golang.org/x/oauth2`).
  - Proxies, middlewares de observabilidad o métricas.
  - Servidores de pruebas `httptest.Server`.
