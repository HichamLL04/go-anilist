# Alcance y Límites de la Librería — `go-anilist`

Este documento define con precisión el perímetro funcional de `go-anilist`.
La librería **no** pretende modelar la totalidad del esquema GraphQL de AniList,
sino cubrir de manera óptima y especializada las necesidades de **aniSync** y
herramientas afines de sincronización de medios.

---

## 1. Principio de Delimitación

El esquema oficial de AniList v2 contiene cientos de tipos, interfaces y uniones
(incluyendo foros, transacciones, feeds sociales, reseñas y actores de doblaje).
Implementar un cliente universal supondría:
- Multiplicar el tamaño del binario y el tiempo de compilación.
- Aumentar las allocations de memoria en runtime sin aportar valor.
- Generar una carga de mantenimiento desproporcionada.

Por ello, `go-anilist` aplica una estrategia de **alcance mínimo viable estricto**,
centrándose en el ciclo de vida de búsqueda, consulta y sincronización de progreso de anime.

---

## 2. Componentes DENTRO del Alcance

| Componente / Servicio | Métodos y Entidades | Justificación técnica en aniSync |
|-----------------------|---------------------|-----------------------------------|
| **Auth / OAuth2** | `AuthorizeURL`, `ExchangeCode`, `RefreshToken` | Permite vincular la cuenta de AniList del usuario de forma desatendida y persistente. |
| **User (`UserService`)** | `Viewer(ctx)`, `GetByID(ctx, id)` | Permite identificar al usuario autenticado, su ID numérico y sus preferencias de lista (sistema de puntuación: 10, 100, 5 estrellas). |
| **Media (`MediaService`)** | `GetByID(ctx, id)`, `GetByIDMal(ctx, idMal)`, `Search(ctx, query, opts...)` | Permite mapear los títulos detectados en Plex (Romaji, Inglés, Nativo) y correlacionar con IDs de MyAnimeList. Incluye conteo de episodios y formato (TV, Movie, OVA, etc.). |
| **MediaList (`MediaListService`)** | `GetEntry(ctx, userID, mediaID)`, `GetCollection(ctx, userID, opts...)`, `SaveEntry(ctx, input)`, `DeleteEntry(ctx, id)` | Núcleo de sincronización: leer si el usuario ya está viendo el anime, su progreso actual, y mutar el episodio visto y estado (`CURRENT`, `COMPLETED`). |
| **Transporte & Cuotas** | `DoGraphQL`, lectura de cabeceras `X-RateLimit-*`, manejo de error 429 | Garantiza no exceder el límite de 90 peticiones/minuto de AniList y reportar cuota a aniSync. |

---

## 3. Componentes FUERA del Alcance

| Área del Esquema AniList | Estado | Justificación de exclusión |
|--------------------------|--------|----------------------------|
| **Foros e Hilos (`Thread`, `Comment`)** | **Excluido** | Irrelevante para sincronización de progreso de vídeo. |
| **Feed Social y Mensajería (`Activity`, `Like`, `Follow`)** | **Excluido** | aniSync es un daemon desatendido; no consume redes sociales. |
| **Reseñas y Recomendaciones (`Review`, `Recommendation`)** | **Excluido** | Alto coste de tipos sin aplicación práctica en el matching. |
| **Personajes, Actores de Voz y Staff (`Character`, `Staff`)** | **Excluido** | Voluminoso y sin uso en el procesamiento de webhooks de Plex. |
| **Notificaciones (`Notification`)** | **Excluido** | No interviene en la sincronización. |
| **Manga y Novelas Ligeras (Lectura/Tracking)** | **Diferido** | aniSync se enfoca en vídeo de Plex. Los enums de soporte base (`MediaType`) se incluyen para no bloquear el futuro, pero no se implementan servicios específicos de manga en esta fase. |
| **Mutaciones administrativas (`UpdateUser`, `Mod*`)** | **Excluido** | La gestión de cuenta y moderación se realiza desde la interfaz web de AniList. |
