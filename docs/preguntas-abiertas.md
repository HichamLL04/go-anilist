# Preguntas Abiertas y Decisiones Pendientes — `go-anilist`

Registro de cuestiones de diseño pendientes de resolución o validación práctica.

---

## 1. Rate Limiter proactivo en la librería vs delegado al consumidor

- **Pregunta**: AniList limita a 90 peticiones por minuto. ¿Debe `go-anilist` incluir un limitador de frecuencia proactivo (ej. `golang.org/x/time/rate`), o debe limitarse a devolver las cabeceras `RateLimitInfo` y gestionar el error `429 Retry-After`, dejando que el consumidor (`aniSync`) controle la tasa de llamadas?
- **Opciones**:
  - *Opción A (Recomendada)*: Mantener `go-anilist` con 0 dependencias externas, capturando pasivamente las cabeceras y exponiendo `ErrRateLimited`. La cola y ritmo de sincronización se controlan en `aniSync`.
  - *Opción B*: Incorporar `golang.org/x/time/rate` opcionalmente activable vía `ClientOption`.
- **Estado**: Pendiente de decisión final al implementar la Fase 1.

---

## 2. Normalización de formatos de puntuación (`ScoreFormat`)

- **Pregunta**: AniList permite al usuario elegir en su perfil diferentes formatos de nota:
  - `POINT_100`: Puntuación de 0 a 100.
  - `POINT_10_DECIMAL`: Puntuación de 0.0 a 10.0 (un decimal).
  - `POINT_10`: Puntuación entera de 0 a 10.
  - `POINT_5`: Escala de 1 a 5 estrellas.
  - `POINT_3`: 3 caras (Smiley / Neutral / Sad).
  ¿Debe `go-anilist` realizar conversiones automáticas o recibir el valor numérico en crudo?
- **Opciones**:
  - *Opción A (Recomendada)*: Recibir el valor crudo en `SaveMediaListEntryInput.Score (*float64)` y dejar que el consumidor adapte según `User.MediaListOptions.ScoreFormat`.
  - *Opción B*: Proveer métodos helper de ayuda (`Score10(val)`, `Score100(val)`).
- **Estado**: Abierta.

---

## 3. Retención opcional del cuerpo HTTP bruto en `Response`

- **Pregunta**: `go-myanimelist` almacena `Response.Body []byte` para que el llamador pueda inspeccionar la respuesta en caso de error. En entornos con memoria restringida, esto duplica los bytes en el heap de Go.
- **Opciones**:
  - *Opción A (Recomendada)*: No retener `Body []byte` en peticiones exitosas (decodificar en streaming). Solo en respuestas con error (`HTTP != 200` o fallos de decode) retener los primeros bytes en `ErrorResponse` para depuración.
  - *Opción B*: Ofrecer una opción de depuración (`WithDebug(true)`) que imprima a `io.Writer`.
- **Estado**: Pre-aprobada Opción A.
