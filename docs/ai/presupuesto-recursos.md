# Presupuesto de Recursos y Memoria — go-anilist

`go-anilist` está diseñada para operar como dependencia central de `aniSync`,
cuyo entorno de ejecución tiene un presupuesto estricto de memoria RAM
(≤ 256 MiB en reposo y ≤ 512 MiB en picos). Por tanto, la librería debe ser
altamente eficiente y predecible.

---

## 1. Coste en tamaño de binario

- **Cero clientes GraphQL sobredimensionados**: No incorporar librerías de terceros
  que arrastren esquemas completos, parsers AST en runtime o introspección pesada.
- **Uso preferente de `net/http` y `encoding/json`**: Mantener las dependencias directas
  en el mínimo absoluto (idealmente sólo biblioteca estándar de Go).
- **Coste medido [RE]**: El impacto de enlazar `go-anilist` en un binario Go final debe
  permanecer por debajo de ~1.5 - 2.0 MiB sin strip y < 800 KiB con `-s -w`.

---

## 2. Gestión de memoria y allocations

- **Decodificación por streaming**: Utilizar `json.NewDecoder(resp.Body)` directamente
  sobre el cuerpo de la respuesta HTTP, cerrando el stream con `defer resp.Body.Close()`.
- **Evitar duplicación de buffers**: No retener copias de respuestas HTTP completas
  en memoria salvo cuando sea indispensable para depuración o reporte de errores.
- **Tipos compactos**: Diseñar structs específicos que contengan únicamente los campos
  solicitados en la consulta GraphQL. Al solicitar sólo los campos necesarios a AniList,
  el payload JSON es significativamente menor en la red y en memoria.

---

## 3. Paginación y control de listas

- **Límite de paginación AniList**: AniList restringe el parámetro `perPage` a un máximo
  de 50 elementos por página **[DOC-ANILIST]**.
- **Prohibido buffering acumulativo forzado**: La librería no debe forzar la carga en
  memoria de toda la lista de medios de un usuario de golpe.
- **Iteración paginada**: Proveer métodos de paginación que permitan al consumidor
  procesar página por página, delegando la agregación en caso de ser requerida.
