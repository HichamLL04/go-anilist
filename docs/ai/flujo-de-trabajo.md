# Flujo de Trabajo para Agentes IA — go-anilist

Instrucciones sobre la metodología y el flujo de trabajo riguroso que todo agente
debe seguir al intervenir en este repositorio.

---

## 1. Principio fundamental: "Leer antes de escribir"

1. **Consultar las fuentes de verdad**:
   - `AGENTS.md` para reglas generales.
   - `docs/reference/anilist-api.md` para validar campos y tipos de la API oficial de AniList.
   - `docs/design/api-publica.md` para verificar firmas y convenciones Go.
2. **No inventar APIs ni campos**:
   - Todo campo añadido a una struct debe coincidir con la documentación oficial de AniList GraphQL v2.
   - Toda firma pública debe seguir las decisiones de diseño preacordadas.

---

## 2. Gestión de Dependencias y Decisiones

- **Prohibido añadir dependencias arbitrarias**: Cualquier adición de módulo en `go.mod`
  debe pasar por un análisis previo en `docs/design/`.
- **Registro de dudas**: Ante cualquier ambigüedad en el comportamiento de AniList o
  del cliente, registrar la consulta en `docs/preguntas-abiertas.md` antes de asumir una solución.

---

## 3. Ciclo de Implementación de una Característica

1. **Definir la consulta GraphQL**: Documentar la query o mutation exacta y los tipos de entrada.
2. **Definir los tipos Go**: Crear las estructuras Go fuertemente tipadas correspondientes.
3. **Escribir el test con fixture**:
   - Crear una respuesta JSON simulada en `testdata/`.
   - Implementar el test unitario con `httptest.Server` validando headers, query enviada y parsing de salida.
4. **Implementar el método del servicio**: Conectar la capa de transporte con el servicio.
5. **Verificar calidad**:
   - Ejecutar `go test -race ./...`.
   - Ejecutar `golangci-lint run`.
6. **Commit atómico**: Realizar commit siguiendo Conventional Commits 1.0.0.
