# Análisis de Referencia: `go-myanimelist`

Este documento analiza en detalle la arquitectura, patrones y decisiones técnicas de
la librería [`go-myanimelist`](https://github.com/nstratos/go-myanimelist)
(autor: Stratos Neiros), utilizada como referencia de diseño para `go-anilist`.

---

## 1. Atribución y Licencia

- **Licencia**: Licencia MIT **[DOC-GO]**.
  ```
  Copyright (c) 2015-present Stratos Neiros <nstratos@gmail.com>
  ```
- **Condición**: La licencia MIT permite reutilizar la arquitectura, patrones y estructuras
  derivadas siempre que se conserve el aviso de copyright y la nota de licencia correspondiente
  en caso de reutilización directa de código.

---

## 2. Estructura de Paquetes

- `go-myanimelist` organiza todo su núcleo en un único paquete principal (`mal`):
  - `mal.go`: Define `Client`, `Response`, `ErrorResponse`, `Do()`, `NewRequest()`, helpers generales.
  - `anime.go`, `manga.go`, `user.go`, `forum.go`: Servicios y modelos de dominio.
  - `useranimelist.go`, `usermangalist.go`: Sub-recursos de listas de usuario.
  - `*_test.go`: Tests unitarios y suite de integración.
- **Evaluación**: Un solo paquete evita fragmentación innecesaria, problemas de importación
  cíclica y jerarquías de paquetes redundantes en bibliotecas cliente Go de tamaño moderado.

---

## 3. Patrón de Cliente y Servicios

El cliente actúa como contenedor central y delega en sub-servicios tipados:

```go
type Client struct {
    client  *http.Client
    BaseURL *url.URL

    Anime *AnimeService
    Manga *MangaService
    User  *UserService
    Forum *ForumService
}

func NewClient(httpClient *http.Client) *Client {
    if httpClient == nil {
        httpClient = &http.Client{}
    }
    baseURL, _ := url.Parse(defaultBaseURL)
    c := &Client{client: httpClient, BaseURL: baseURL}
    c.Anime = &AnimeService{client: c}
    c.Manga = &MangaService{client: c}
    c.User  = &UserService{client: c}
    return c
}
```

Cada servicio almacena una referencia no exportada a `client *Client`, lo que le permite
acceder al método `c.Do(ctx, req, v)` y al transporte HTTP sin duplicar lógica de red.

---

## 4. Opciones Funcionales Tipadas

En lugar de listas largas de parámetros o structs de configuración repletas de punteros nil,
`go-myanimelist` utiliza opciones funcionales basadas en interfaces:

```go
type Option interface {
    apply(v *url.Values)
}

type Limit int
func (l Limit) apply(v *url.Values) { v.Set("limit", strconv.Itoa(int(l))) }

type Fields []string
func (f Fields) apply(v *url.Values) {
    if len(f) != 0 {
        v.Set("fields", strings.Join(f, ","))
    }
}
```

Para métodos específicos, crea interfaces especializadas (`AnimeListOption`, `UpdateMyAnimeListStatusOption`),
lo que impide en tiempo de compilación pasar opciones no válidas a un método concreto.

---

## 5. Autenticación

- `go-myanimelist` no reinventa la rueda de OAuth2: delega la autenticación en
  `golang.org/x/oauth2`.
- El consumidor simplemente inicializa el cliente inyectando el `*http.Client` autenticado
  proporcionado por `oauth2.Config.Client(ctx, token)`.
- Esto separa limpiamente el ciclo de vida de los tokens (almacenamiento, refresh) de la lógica
  propia del SDK.

---

## 6. Paginación

- Al ser una API REST, MyAnimeList devuelve enlaces absolutos en un objeto `Paging`:
  ```json
  "paging": {
    "next": "https://api.myanimelist.net/v2/users/@me/animelist?offset=10&limit=10",
    "previous": "..."
  }
  ```
- `go-myanimelist` parsea el query parameter `offset` de esas URLs y expone `resp.NextOffset`
  y `resp.PrevOffset` en el struct `Response`.

---

## 7. Manejo de Errores

- Define el struct `ErrorResponse`:
  ```go
  type ErrorResponse struct {
      Response *http.Response
      Message  string `json:"message"`
      Err      string `json:"error"`
  }
  ```
- En `checkResponse`, si el código HTTP no está en el rango `200..299`, lee el cuerpo,
  intenta deserializar el error JSON y devuelve `*ErrorResponse`.
- Permite al usuario comprobar errores específicos con `errors.As(err, &errResp)`.

---

## 8. Estrategia de Testing

- Emplea `net/http/httptest.Server` para pruebas totalmente aisladas sin red:
  ```go
  func setup() (client *Client, mux *http.ServeMux, teardown func()) {
      mux = http.NewServeMux()
      server := httptest.NewServer(mux)
      client = NewClient(nil)
      client.BaseURL, _ = url.Parse(server.URL + "/")
      return client, mux, server.Close
  }
  ```
- Funciones auxiliares de verificación: `testMethod`, `testURLValues`, `testBody`, `testErrorResponse`.
- Respuestas de prueba almacenadas o inyectadas como strings JSON mock.

---

## 9. Comparativa: Qué adoptar vs Qué NO copiar para `go-anilist`

| Patrón en `go-myanimelist` | ¿Adoptar en `go-anilist`? | Justificación técnica |
|-----------------------------|---------------------------|-----------------------|
| **Patrón Client / Services** | **SÍ** | Estructura modular idiomática en Go; agrupa métodos por entidad (`Media`, `MediaList`, `Viewer`). |
| **Inyección de `*http.Client`** | **SÍ** | Permite inyectar autenticación OAuth2, rate limiters, proxies y servidores `httptest`. |
| **Opciones funcionales tipadas** | **SÍ (adaptado)** | Excelente ergonomía; en vez de modificar `url.Values`, las opciones modificarán un mapa de `variables` GraphQL. |
| **Struct `Response` enriquecido** | **SÍ** | Fundamental para exponer cabeceras `X-RateLimit-*` y cuotas de AniList al llamador. |
| **Estrategia `httptest.Server`** | **SÍ (adaptado)** | Testing determinista sin dependencias externas ni conexión a internet. |
| **Construcción de URLs REST (`url.Values`)** | **NO** | AniList es 100% GraphQL; todo tráfico va vía HTTP POST con JSON al endpoint raíz. |
| **Paginación por offset en URLs** | **NO** | AniList utiliza GraphQL `Page(page: $page, perPage: $perPage)` y el objeto `PageInfo`. |
| **Detección de errores sólo por HTTP Status** | **NO** | GraphQL puede devolver HTTP 200 con un array `"errors": [...]` en el cuerpo. El parser debe revisar siempre `errors`. |
| **Muxing por URL path en tests** | **NO** | En AniList todas las peticiones apuntan a la misma ruta `/`; el handler de test debe inspeccionar `operationName` o variables en el body JSON. |
