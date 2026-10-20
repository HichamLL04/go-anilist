# Especificación de la API GraphQL de AniList (v2)

Documentación de referencia técnica sobre la API de AniList basada en
la documentación oficial (https://docs.anilist.co) **[DOC-ANILIST]**.

---

## 1. Endpoint y Protocolo HTTP

- **Endpoint de producción**: `https://graphql.anilist.co` **[DOC-ANILIST]**
- **Método**: Exclusivamente `POST`. Peticiones `GET` son rechazadas.
- **Cabeceras obligatorias**:
  ```http
  Content-Type: application/json
  Accept: application/json
  ```
- **Cabecera de autenticación** (cuando se requieren permisos de usuario):
  ```http
  Authorization: Bearer <access_token>
  ```
- **Cuerpo de la petición**:
  ```json
  {
    "query": "query ($id: Int) { Media(id: $id) { id title { romaji english } } }",
    "variables": { "id": 1 }
  }
  ```

---

## 2. Autenticación y OAuth2

AniList utiliza el protocolo OAuth 2.0 para autorizar el acceso y mutaciones en nombre del usuario **[DOC-ANILIST]**:

- **Authorize URL**: `https://anilist.co/api/v2/oauth/authorize`
- **Token URL**: `https://anilist.co/api/v2/oauth/token`

### Flujos soportados

1. **Authorization Code Grant** (Recomendado para servicios backend como aniSync):
   - Redirección del usuario a `/oauth/authorize?client_id={CLIENT_ID}&response_type=code&redirect_uri={REDIRECT_URI}`.
   - Intercambio del `code` mediante POST a `/oauth/token` con payload:
     ```json
     {
       "grant_type": "authorization_code",
       "client_id": "{CLIENT_ID}",
       "client_secret": "{CLIENT_SECRET}",
       "redirect_uri": "{REDIRECT_URI}",
       "code": "{AUTHORIZATION_CODE}"
     }
     ```
   - Devuelve `access_token`, `token_type: "Bearer"`, `expires_in` y `refresh_token`.
2. **Implicit Grant** (Para tokens estáticos o entornos de un solo usuario):
   - Redirección con `response_type=token`.
   - Devuelve el token directamente en el fragmento hash de la URL (`#access_token=...`).

### Caducidad y validez de tokens

- Los tokens de acceso de AniList tienen típicamente una validez de **1 año (31.536.000 segundos)** **[DOC-ANILIST]**.
- Cuando se utiliza Authorization Code Grant, el `refresh_token` permite obtener un nuevo `access_token` antes o después de la expiración sin requerir intervención manual del usuario.

---

## 3. Rate Limiting (Límites de Cuota Reales)

- **Límite general**: **90 peticiones por minuto** por IP o por token OAuth2 **[DOC-ANILIST]**.
- Si se supera el límite, la API responde con código HTTP **`429 Too Many Requests`**.
- **Cabeceras de respuesta HTTP**:
  - `X-RateLimit-Limit`: Límite máximo de peticiones por minuto (típicamente `90`).
  - `X-RateLimit-Remaining`: Peticiones disponibles en la ventana actual de 60 segundos.
  - `X-RateLimit-Reset`: Timestamp UNIX (en segundos) en el que se reinicia el contador de cuota.
  - `Retry-After`: En respuestas HTTP 429, indica los segundos enteros que el cliente debe esperar antes de reintentar.

---

## 4. Gestión de Errores y Formato de Respuesta

AniList devuelve respuestas estructuradas estándar de GraphQL **[DOC-ANILIST]**:

### Respuesta exitosa
```json
{
  "data": { ... }
}
```

### Respuesta con errores
```json
{
  "errors": [
    {
      "message": "Validation failed: Score must be between 0 and 100",
      "status": 400,
      "locations": [{ "line": 2, "column": 3 }]
    }
  ]
}
```

### Códigos de estado HTTP devueltos por AniList
- **200 OK**: Consulta ejecutada (incluso cuando ciertos campos parciales contengan errores en el array `errors`).
- **400 Bad Request**: Error de validación de argumentos, sintaxis GraphQL inválida o violaciones de tipos.
- **401 Unauthorized**: Token de acceso ausente, inválido o caducado.
- **404 Not Found**: Recurso no encontrado (ej. `Media` con ID inexistente).
- **429 Too Many Requests**: Límite de 90 peticiones por minuto excedido.
- **500 / 502 / 503 Server Error**: Incidencia en los servidores de AniList o Cloudflare.

---

## 5. Paginación en AniList

- Toda consulta de colecciones utiliza el envoltorio `Page(page: $page, perPage: $perPage)` **[DOC-ANILIST]**.
- **Restricción estricta**: El valor máximo de `perPage` permitido por la API es **50**. Solicitar más de 50 produce un error o un truncamiento automático a 50.
- **Estructura de `pageInfo`**:
  ```graphql
  pageInfo {
    total
    perPage
    currentPage
    lastPage
    hasNextPage
  }
  ```

---

## 6. Queries y Mutations Requeridas para aniSync

Para cumplir los requisitos de sincronización de aniSync, `go-anilist` debe cubrir exactamente las siguientes operaciones:

### 1. Obtener usuario autenticado (`Viewer`)
Permite conocer el `id`, nombre de usuario y opciones de listas del token configurado.

```graphql
query CurrentUser {
  Viewer {
    id
    name
    avatar {
      medium
    }
    mediaListOptions {
      scoreFormat
    }
  }
}
```

### 2. Búsqueda y consulta de Anime (`Media`)
Permite identificar el anime a partir de títulos detectados en Plex (Romaji, Inglés, Nativo) o IDs cruzados (MyAnimeList).

```graphql
query SearchAnime($search: String, $idMal: Int, $type: MediaType, $format: MediaFormat) {
  Page(page: 1, perPage: 10) {
    media(search: $search, idMal: $idMal, type: $type, format: $format) {
      id
      idMal
      title {
        romaji
        english
        native
      }
      format
      status
      episodes
      duration
      synonyms
      startDate {
        year
        month
        day
      }
    }
  }
}
```

### 3. Consulta de entrada en la lista (`MediaList`)
Verifica el estado actual de seguimiento de un anime concreto para el usuario autenticado.

```graphql
query GetMediaListEntry($userId: Int, $mediaId: Int) {
  MediaList(userId: $userId, mediaId: $mediaId, type: ANIME) {
    id
    mediaId
    status
    progress
    score
    repeat
    updatedAt
  }
}
```

### 4. Colección de listas de usuario (`MediaListCollection`)
Obtiene la lista completa de seguimiento de un usuario de forma paginada/chunked para sincronización o verificación de biblioteca.

```graphql
query GetUserAnimeList($userId: Int, $status: MediaListStatus) {
  MediaListCollection(userId: $userId, type: ANIME, status: $status) {
    lists {
      name
      isCustomList
      status
      entries {
        id
        mediaId
        status
        progress
        score
        repeat
        updatedAt
        media {
          id
          idMal
          title {
            romaji
            english
          }
          episodes
        }
      }
    }
  }
}
```

### 5. Actualización / Creación de entrada (`SaveMediaListEntry`)
Mutación principal para reportar progreso de visualización de episodios y estado (Watching, Completed).

```graphql
mutation UpdateProgress(
  $mediaId: Int
  $status: MediaListStatus
  $progress: Int
  $score: Float
  $repeat: Int
) {
  SaveMediaListEntry(
    mediaId: $mediaId
    status: $status
    progress: $progress
    score: $score
    repeat: $repeat
  ) {
    id
    mediaId
    status
    progress
    score
    repeat
    updatedAt
  }
}
```

### 6. Eliminación de entrada (`DeleteMediaListEntry`)
```graphql
mutation DeleteEntry($id: Int) {
  DeleteMediaListEntry(id: $id) {
    deleted
  }
}
```
