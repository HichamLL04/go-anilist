# Convención de Commits — go-anilist

Todos los commits de este repositorio deben cumplir estrictamente la especificación
**Conventional Commits 1.0.0**.

---

## 1. Estructura del commit

```
<type>(<scope>)<!>: <description>

[cuerpo opcional: líneas de ≤ 72 caracteres]

[pie(s) opcional(es): BREAKING CHANGE:, Refs:]
```

### Reglas obligatorias:
- **type**: Uno de los tipos permitidos (en minúsculas).
- **scope**: Obligatorio siempre que el cambio afecte a un subsistema concreto.
- **description**: En inglés, modo imperativo ("add", "fix", "update", nunca "added" o "adds"), minúscula inicial, sin punto final, longitud máxima 72 caracteres.
- **breaking change**: Indicar `!` inmediatamente tras el tipo/scope (ej. `feat(client)!: change auth signature`) y añadir `BREAKING CHANGE: <desc>` en el footer.
- **atomicidad**: Un commit por cambio lógico. No acumular cambios no relacionados en un único commit masivo.

---

## 2. Tipos permitidos

| Tipo | Propósito | Ejemplo de uso |
|------|-----------|----------------|
| `feat` | Nueva característica o función pública | Implementar servicio de MediaList |
| `fix` | Corrección de un bug o comportamiento anómalo | Corregir parsing de errores GraphQL |
| `docs` | Cambios exclusivos en documentación | Añadir especificación de API |
| `style` | Formato, espacios, punto y coma (sin cambios de código) | Ajustar formato gofmt |
| `refactor` | Refactorización de código sin cambiar comportamiento ni API | Simplificar decodificador interno |
| `perf` | Optimización de rendimiento o asignaciones de memoria | Reutilizar buffers en JSON decoder |
| `test` | Añadir o corregir tests unitarios | Añadir fixtures httptest para búsqueda |
| `build` | Cambios en sistema de build o dependencias de tooling | Ajustar configuración de Go |
| `ci` | Cambios en flujos de integración continua | Añadir workflow de GitHub Actions |
| `chore` | Tareas de mantenimiento o configuración de repo | Inicializar gitignore o licencias |
| `revert` | Reversión de un commit anterior | Revertir commit previo |

---

## 3. Scopes comunes en go-anilist

- `client`: Núcleo del cliente HTTP / GraphQL base.
- `media`: Servicio de consulta y búsqueda de Media (Anime/Manga).
- `medialist`: Servicio de consulta y actualización de listas de usuario.
- `user`: Servicio de perfil de usuario y `Viewer`.
- `auth`: Autenticación OAuth2 y gestión de tokens.
- `graphql`: Motor de transporte GraphQL, generación de payload y parsing de errores.
- `ratelimit`: Manejo de headers de cuota (`X-RateLimit-*`) y backoff.
- `design`: Documentos de diseño de arquitectura.
- `reference`: Documentos de análisis comparativo y especificaciones.
- `ai`: Reglas y configuraciones para agentes IA.
- `repo`: Configuración general del repositorio.

---

## 4. Ejemplos de commits

### Buenos ejemplos

```text
docs(reference): document go-myanimelist design analysis and patterns

Analyze package structure, service delegation, functional options,
and testing techniques for comparison with AniList GraphQL needs.
```

```text
feat(medialist): add UpdateMediaListEntry method and options

Implement SaveMediaListEntry mutation wrapper supporting score,
status, progress and repeat increments.
```

```text
fix(graphql): handle GraphQL error arrays with missing HTTP status

AniList returns HTTP 200 with top-level "errors" array when schema
validation or field resolution fails.
```

### Malos ejemplos (PROHIBIDOS)

- ❌ `added some files`: No usa Conventional Commits, descripción en pasado.
- ❌ `WIP`: Commits parciales sin mensaje descriptivo.
- ❌ `docs: changes.`: Falta scope cuando aplica, contiene punto final.
- ❌ `feat(client): Add new method and also update readme and fix bug`: No es atómico; mezcla 3 intenciones distintas.
- ❌ `feat(client): implement huge client refactoring with all services`: Commit gigante monolítico.
