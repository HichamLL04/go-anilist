# Estrategia de Testing — `go-anilist`

Este documento describe la arquitectura y estándares de pruebas para `go-anilist`,
asegurando suites deterministas, rápidas, aisladas y reproducibles sin conexión a red.

---

## 1. Principios Fundamentales

1. **Cero tráfico de red real**: Ningún test unitario debe conectarse a `https://graphql.anilist.co`. Todo tráfico se simula localmente.
2. **Determinismo y velocidad**: La suite completa debe completarse en menos de 2 segundos en cualquier máquina de desarrollo o pipeline de CI.
3. **Concurrencia segura**: Toda ejecución de tests debe superar el detector de condiciones de carrera de Go (`go test -race ./...`).
4. **Validación bidireccional**: Los tests deben verificar tanto lo que el cliente *envía* (query GraphQL, variables, headers HTTP) como lo que el cliente *recibe y parsea* (tipos Go, errores enriquecidos, cuotas).

---

## 2. Servidor de Prueba Local (`httptest.Server`)

A diferencia de `go-myanimelist` (donde cada recurso tiene una URL REST distinta como `/anime/1` o `/users/@me`),
en AniList GraphQL todas las peticiones se dirigen mediante `POST` a la raíz `/`.

Por ello, el multiplexor de tests debe inspeccionar el payload JSON:

```go
type testMuxHandler struct {
    handlers map[string]http.HandlerFunc // Clave: operationName o query match
}

func setupTestServer(t *testing.T) (*Client, *testMuxHandler, func()) {
    t.Helper()
    muxHandler := &testMuxHandler{
        handlers: make(map[string]http.HandlerFunc),
    }

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Validar método y cabeceras obligatorias
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        if r.Header.Get("Content-Type") != "application/json" {
            http.Error(w, "bad content type", http.StatusBadRequest)
            return
        }

        // 2. Decodificar la petición GraphQL para saber qué operación se ejecuta
        var reqBody struct {
            Query         string                 `json:"query"`
            OperationName string                 `json:"operationName"`
            Variables     map[string]interface{} `json:"variables"`
        }
        if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
            http.Error(w, "invalid json body", http.StatusBadRequest)
            return
        }

        // 3. Despachar al handler específico registrado
        opKey := reqBody.OperationName
        if opKey == "" {
            opKey = detectOperationKey(reqBody.Query)
        }

        if handler, ok := muxHandler.handlers[opKey]; ok {
            handler(w, r)
            return
        }

        http.Error(w, fmt.Sprintf("unhandled test operation: %s", opKey), http.StatusNotFound)
    }))

    client := NewClient(nil, WithBaseURL(server.URL))
    return client, muxHandler, server.Close
}
```

---

## 3. Fixtures y Datos de Prueba (`testdata/`)

Las respuestas simuladas se organizan en ficheros JSON estáticos dentro de `testdata/`:

```
testdata/
  viewer_success.json              ← Respuesta a query Viewer
  media_search_success.json        ← Búsqueda con 3 resultados y PageInfo
  medialist_get_success.json       ← Entrada existente en lista
  medialist_save_success.json      ← Mutación SaveMediaListEntry completada
  error_graphql_validation.json    ← Array "errors" con status 400
  error_rate_limited.json          ← Respuesta 429 con cabecera Retry-After
```

Se cargan en los tests mediante `os.ReadFile("testdata/...")` o variables string precompiladas.

---

## 4. Estructura de Tests Table-Driven

Todos los casos de prueba deben estructurarse en tablas `tc`:

```go
func TestMediaService_Search(t *testing.T) {
    cases := []struct {
        name          string
        query         string
        opts          []SearchOption
        mockStatus    int
        mockHeaders   map[string]string
        mockFixture   string
        wantErr       error
        wantCount     int
        wantRemaining int
    }{
        {
            name:          "successful anime search",
            query:         "Frieren",
            opts:          []SearchOption{WithFormat(MediaFormatTV)},
            mockStatus:    http.StatusOK,
            mockHeaders:   map[string]string{"X-RateLimit-Remaining": "85"},
            mockFixture:   "testdata/media_search_success.json",
            wantErr:       nil,
            wantCount:     1,
            wantRemaining: 85,
        },
        {
            name:        "rate limited 429 response",
            query:      "Attack on Titan",
            mockStatus: http.StatusTooManyRequests,
            mockHeaders: map[string]string{
                "Retry-After": "45",
            },
            mockFixture: "testdata/error_rate_limited.json",
            wantErr:     ErrRateLimited,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            // Setup server, configure mock handler, execute method and assert
        })
    }
}
```

---

## 5. Verificación de Cancelación de Contexto

Cada servicio debe incluir tests explícitos que garanticen que un contexto cancelado (`context.Canceled` o `context.DeadlineExceeded`) interrumpe la petición de inmediato sin fugas de goroutines ni recursos de socket.
