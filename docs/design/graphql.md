# Estrategia de Integración GraphQL — `go-anilist`

Este documento analiza y compara tres enfoques técnicos para la interacción con
la API GraphQL de AniList en Go, evaluando dependencias, coste de binario,
mantenimiento y asignaciones de memoria, y formula una recomendación formal.

---

## 1. Opciones Técnicas Evaluadas

### Opción A: Consultas manuales con structs tipados y `net/http` + `encoding/json`
- **Descripción**: Las queries y mutations se definen como constantes string en Go (`const searchAnimeQuery = "..."`). Las respuestas se deserializan directamente en structs Go tipados usando `json.NewDecoder` de la biblioteca estándar.
- **Dependencias externas**: **0**. Utiliza exclusivamente la biblioteca estándar de Go (`net/http`, `encoding/json`, `context`, `fmt`).

### Opción B: Generación de código tipado (ej. `genqlient`)
- **Descripción**: Se descarga el archivo de esquema SDL completo de AniList (`schema.graphql`). Se configuran queries en archivos `.graphql` y `genqlient` genera código Go fuertemente tipado en tiempo de compilación.
- **Dependencias externas**: Requiere `github.com/Khan/genqlient` como herramienta de build (`tools.go`), configuración `genqlient.yaml` y pasos de `go generate`.

### Opción C: Clientes GraphQL genéricos dinámicos (ej. `machinebox/graphql`, `shurcooL/graphql`)
- **Descripción**: Bibliotecas de terceros que encapsulan peticiones GraphQL mediante reflexión en tiempo de ejecución o abstracciones sobre `http.Client`.
- **Dependencias externas**: 1 o más librerías externas en `go.mod`.

---

## 2. Matriz Comparativa

| Criterio | Opción A: Consultas Manuales | Opción B: Generador (`genqlient`) | Opción C: Clientes Genéricos |
|----------|------------------------------|-----------------------------------|------------------------------|
| **Dependencias en `go.mod`** | **0 (Solo stdlib)** **[DOC-GO]** | 1 dependencia de tooling (`genqlient`) | 1+ dependencias externas |
| **Coste en binario [RE]** | **Mínimo absoluto (~30-50 KiB de código compilado)** | Moderado (~250-400 KiB en structs generados) | Medio (~300-600 KiB por reflexión) |
| **Paso de compilación** | Inmediato (`go build ./...`) | Requiere `go generate ./...` y esquema SDL | Inmediato (`go build ./...`) |
| **Mantenimiento del esquema** | Solo se definen los campos exactos que se consumen | Requiere mantener sincronizado un SDL de ~2-3 MB | Mantenimiento dinámico manual |
| **Uso de memoria (allocs)** | Óptimo (streaming directo con `json.NewDecoder`) | Muy bueno (structs directos) | Regular (overhead de reflexión en runtime) |
| **Control de peticiones de red** | Control total byte a byte de payload y cabeceras | Control indirecto según generador | Abstracción rígida sobre HTTP |

---

## 3. Análisis Detallado

### Coste en binario y dependencias
En `aniSync`, el presupuesto de RAM es crítico (≤ 256 MiB en reposo y ≤ 512 MiB en picos). Añadir dependencias de terceros para ejecutar apenas 5 queries y 2 mutations introduce riesgos de bloating y vulnerabilidades en la cadena de suministro. La Opción A no agrega ninguna dependencia a `go.mod`.

### Superficie de mantenimiento
AniList añade periódicamente campos y mutaciones a su esquema general. Con `genqlient` (Opción B), actualizar el esquema requiere procesar un archivo SDL masivo que incluye cientos de tipos irrelevantes para aniSync (como sistemas de foros y moderación). En cambio, con la Opción A, cada consulta solicita exclusivamente los campos necesarios, garantizando payloads JSON pequeños y tiempos de deserialización mínimos.

### Control sobre cabeceras de rate limiting
AniList expone las métricas de consumo de API mediante cabeceras HTTP (`X-RateLimit-*`). Muchas librerías genéricas de GraphQL (Opción C) aíslan tanto el transporte HTTP que resulta complejo o ineficiente interceptar cabeceras de respuesta por petición sin recurrir a `RoundTripper` personalizados complejos. Con la Opción A, el control sobre el objeto `*http.Response` es directo y transparente.

---

## 4. Decisión y Recomendación

**Se adopta la Opción A: Consultas manuales con structs tipados de la biblioteca estándar.**

### Justificación:
1. **Cero dependencias externas**: Cumple estrictamente con la política de estabilidad y ligereza.
2. **Máxima eficiencia**: `json.NewDecoder` decodifica por streaming evitando duplicar buffers en memoria.
3. **Simplicidad**: El alcance de `go-anilist` consta de exactamente 4 consultas y 2 mutaciones. No se justifica la complejidad de un pipeline de generación de código SDL.
4. **Inspección de cuotas**: Permite capturar las cabeceras `X-RateLimit-*` y responder a códigos `429` de manera nativa.
