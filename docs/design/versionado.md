# Política de Versionado y Ciclo de Vida — `go-anilist`

Este documento establece el esquema de versionado, la relación con
Conventional Commits y la política de compatibilidad hacia atrás de `go-anilist`.

---

## 1. Esquema SemVer 2.0.0

`go-anilist` sigue estrictamente **Semantic Versioning 2.0.0** (`vMAJOR.MINOR.PATCH`) **[DOC-GO]**:

- **MAJOR**: Cambios incompatibles en la API pública (breaking changes).
- **MINOR**: Nuevas funcionalidades públicas compatibles hacia atrás (ej. nuevos campos o métodos de servicio).
- **PATCH**: Corrección de errores compatibles hacia atrás o mejoras internas de rendimiento.

---

## 2. Fase Inicial: Ciclo `v0.x`

1. **Punto de partida**: La primera versión funcional etiquetada será **`v0.1.0`**.
2. **Evolución en `v0.x`**:
   - De acuerdo con la especificación SemVer, durante la serie `v0.x` la API pública puede evolucionar rápidamente para adaptarse a las necesidades de integración con `aniSync`.
   - Los cambios que rompan compatibilidad en esta fase incrementarán el dígito `MINOR` (ej. `v0.1.0` -> `v0.2.0`).
   - Las adiciones y correcciones menores incrementarán el dígito `PATCH` (ej. `v0.1.0` -> `v0.1.1`).
3. **Transición a `v1.0.0`**:
   - `go-anilist` se promoverá a `v1.0.0` únicamente tras validar en producción que la suite de `aniSync` opera de manera estable durante un ciclo completo y las firmas de la API pública han quedado consolidadas.

---

## 3. Relación entre Tipos de Commit y Versiones

La automatización de releases deriva directamente del historial de Git bajo Conventional Commits:

| Tipo de Commit | Efecto en `v0.x` | Efecto en `v1.x+` |
|----------------|------------------|-------------------|
| `fix(...)` | `PATCH` (`v0.1.0` -> `v0.1.1`) | `PATCH` (`v1.0.0` -> `v1.0.1`) |
| `perf(...)` | `PATCH` (`v0.1.0` -> `v0.1.1`) | `PATCH` (`v1.0.0` -> `v1.0.1`) |
| `feat(...)` | `PATCH` o `MINOR` según entidad | `MINOR` (`v1.0.0` -> `v1.1.0`) |
| `feat(...)!` o `BREAKING CHANGE:` | `MINOR` (`v0.1.0` -> `v0.2.0`) | `MAJOR` (`v1.0.0` -> `v2.0.0`) |
| `docs(...)`, `style(...)`, `chore(...)`, `test(...)` | Sin incremento de versión | Sin incremento de versión |

---

## 4. Etiquetado y Módulos Go

- Los tags de Git deben llevar el prefijo `v` (ej. `v0.1.0`, `v0.1.1`).
- En `go.mod`, la ruta del módulo se mantendrá canónica:
  ```go
  module github.com/morluto/go-anilist // o la organización final
  ```
- Al avanzar a versiones mayores (`v2+`), se adoptará la convención de sufijo de módulo de Go (`/v2`).
