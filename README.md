# go-anilist

[![Go Reference](https://pkg.go.dev/badge/github.com/hichamll04/go-anilist.svg)](https://pkg.go.dev/github.com/hichamll04/go-anilist)
[![CI](https://github.com/hichamll04/go-anilist/actions/workflows/ci.yml/badge.svg)](https://github.com/hichamll04/go-anilist/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Librería cliente Go ligera e independiente para la API GraphQL (v2) de [AniList](https://anilist.co).
Diseñada con **cero dependencias externas** (solo biblioteca estándar de Go), decodificación por streaming y estricta eficiencia de memoria para aplicaciones de sincronización y daemons.

---

## Características

- **Zero External Dependencies**: Desarrollada exclusivamente sobre la biblioteca estándar de Go (`net/http`, `encoding/json`, `context`).
- **Streaming JSON**: Deserialización directa sin buffering acumulativo redundante.
- **Rate Limit Aware**: Inspección automática de cabeceras `X-RateLimit-*` y gestión nativa de respuestas `429 Too Many Requests` (`Retry-After`).
- **Servicios Modulares**:
  - `User`: Perfil autenticado (`Viewer`) y configuraciones de listas (`ScoreFormat`).
  - `Media`: Búsqueda de anime y resolución por ID AniList / ID MyAnimeList.
  - `MediaList`: Consulta de progreso de visualización, colecciones y actualización (`SaveEntry`).
  - `Auth`: Flujo de autorización OAuth2 (`AuthorizeURL`, `ExchangeCode`, `RefreshToken`).
- **100% Determinista en Tests**: Suite de pruebas sin tráfico de red real mediante `net/http/httptest`.

---

## Instalación

```bash
go get github.com/hichamll04/go-anilist
```

---

## Ejemplos de Uso

### 1. Inicialización del Cliente

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/hichamll04/go-anilist"
)

func main() {
    // Inicializar cliente con token Bearer
    client := anilist.NewClient(nil, anilist.WithToken("tu_access_token_aqui"))

    // Obtener información del usuario autenticado
    user, resp, err := client.User.CurrentUser(context.Background())
    if err != nil {
        log.Fatalf("Error obteniendo usuario: %v", err)
    }

    fmt.Printf("Usuario: %s (ID: %d)\n", user.Name, user.ID)
    fmt.Printf("Cuota restante: %d peticiones\n", resp.RateLimit.Remaining)
}
```

### 2. Búsqueda de Anime

```go
results, _, err := client.Media.Search(
    ctx,
    "Sousou no Frieren",
    anilist.WithFormat(anilist.MediaFormatTV),
    anilist.WithSeasonYear(2023),
    anilist.WithLimit(5),
)
if err != nil {
    log.Fatal(err)
}

for _, anime := range results {
    fmt.Printf("[%d] %s - %d episodios\n", anime.ID, anime.Title.Romaji, *anime.Episodes)
}
```

### 3. Actualizar Progreso de Visualización

```go
progress := 16
status := anilist.MediaListStatusCurrent
score := 9.5

entry, _, err := client.MediaList.SaveEntry(ctx, anilist.SaveMediaListEntryInput{
    MediaID:  154587,
    Progress: &progress,
    Status:   &status,
    Score:    &score,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Actualizado anime %d al episodio %d (estado: %s)\n", entry.MediaID, entry.Progress, entry.Status)
```

### 4. Manejo Estructurado de Errores

```go
if errors.Is(err, anilist.ErrNotFound) {
    fmt.Println("El recurso no existe en AniList")
} else if errors.Is(err, anilist.ErrRateLimited) {
    fmt.Println("Límite de peticiones excedido (429)")
} else if errors.Is(err, anilist.ErrUnauthorized) {
    fmt.Println("Token inválido o expirado")
}
```

---

## Agradecimientos y Atribución

La arquitectura de este cliente (patrón de sub-servicios `Client.Service`, inyección limpia de `*http.Client`, opciones funcionales y estrategia de testing determinista con `httptest.Server`) está inspirada en la excelente librería [`go-myanimelist`](https://github.com/nstratos/go-myanimelist) desarrollada por **Stratos Neiros** (licencia MIT). 

`go-anilist` adapta dichos principios al protocolo GraphQL v2 de AniList, manteniendo total compatibilidad idiomática en Go y cero dependencias externas.

---

## Licencia

Distribuido bajo la [Licencia MIT](LICENSE).